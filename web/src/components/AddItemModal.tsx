import React, { useState } from 'react';
import { Bookmark, Inbox, Loader2, Plus, Sparkles, X } from 'lucide-react';
import { ItemType } from '../types';
import { api } from '../api/client';

interface AddItemModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (itemData: {
    type: ItemType;
    title: string;
    creator: string;
    external_url: string;
    cover_url?: string;
    destination: 'queue' | 'shelf';
    rating?: number;
    take?: string;
  }) => Promise<void>;
}

const MEDIA_TYPES: { type: ItemType; label: string; icon: string }[] = [
  { type: 'movie', label: 'Movie / TV', icon: '🎬' },
  { type: 'youtube', label: 'YouTube', icon: '▶️' },
  { type: 'book', label: 'Book', icon: '📖' },
  { type: 'anime', label: 'Anime', icon: '✨' },
  { type: 'podcast', label: 'Podcast', icon: '🎧' },
];

export const AddItemModal: React.FC<AddItemModalProps> = ({ isOpen, onClose, onSubmit }) => {
  const [type, setType] = useState<ItemType>('movie');
  const [title, setTitle] = useState('');
  const [creator, setCreator] = useState('');
  const [externalUrl, setExternalUrl] = useState('');
  const [coverUrl, setCoverUrl] = useState('');
  const [destination, setDestination] = useState<'queue' | 'shelf'>('queue');
  const [rating, setRating] = useState<number>(8);
  const [take, setTake] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isEnriching, setIsEnriching] = useState(false);
  const [enrichSuccess, setEnrichSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  // Auto-detect type & fetch rich metadata via Phase 2 pipeline
  const handleUrlChange = async (url: string) => {
    setExternalUrl(url);
    setEnrichSuccess(false);

    const lower = url.toLowerCase();
    if (lower.includes('youtube.com') || lower.includes('youtu.be')) {
      setType('youtube');
    } else if (lower.includes('themoviedb.org') || lower.includes('imdb.com') || lower.includes('letterboxd.com')) {
      setType('movie');
    } else if (lower.includes('openlibrary.org') || lower.includes('goodreads.com')) {
      setType('book');
    } else if (lower.includes('anilist.co') || lower.includes('myanimelist.net')) {
      setType('anime');
    } else if (lower.includes('spotify.com') || lower.includes('apple.com/podcast')) {
      setType('podcast');
    }

    // If valid URL with protocol, attempt automated metadata enrichment
    if (url.startsWith('http://') || url.startsWith('https://')) {
      try {
        setIsEnriching(true);
        const res = await api.enrichURL(url);
        if (res.success && res.data) {
          const meta = res.data;
          if (meta.title && !title) setTitle(meta.title);
          if (meta.creator && !creator) setCreator(meta.creator);
          if (meta.cover_url && !coverUrl) setCoverUrl(meta.cover_url);
          if (meta.type) setType(meta.type);
          setEnrichSuccess(true);
        }
      } catch {
        // Silent fallback to manual entry
      } finally {
        setIsEnriching(false);
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      setError('Title is required');
      return;
    }
    if (!externalUrl.trim()) {
      setError('External URL is required');
      return;
    }

    try {
      setIsSubmitting(true);
      setError(null);
      await onSubmit({
        type,
        title: title.trim(),
        creator: creator.trim(),
        external_url: externalUrl.trim(),
        cover_url: coverUrl.trim() || undefined,
        destination,
        rating: destination === 'shelf' ? rating : undefined,
        take: destination === 'shelf' ? take.trim() : undefined,
      });
      // Reset form
      setTitle('');
      setCreator('');
      setExternalUrl('');
      setCoverUrl('');
      setTake('');
      onClose();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to save item');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-sand-950/40 backdrop-blur-xs">
      <div
        className="bg-white rounded-2xl max-w-lg w-full border border-sand-200 shadow-xl overflow-hidden animate-in fade-in zoom-in-95 duration-150"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-sand-100 bg-sand-50/50">
          <div>
            <h3 className="text-base font-serif font-semibold text-sand-900">
              Save Recommendation
            </h3>
            <p className="text-xs text-sand-500">
              Add culture to your personal queue or taste journal.
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-lg text-sand-400 hover:text-sand-700 hover:bg-sand-100 transition"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          {/* Destination Selector: Queue vs Shelf */}
          <div>
            <label className="block text-xs font-semibold text-sand-700 uppercase tracking-wider mb-2">
              Destination
            </label>
            <div className="grid grid-cols-2 gap-2 p-1 bg-sand-100 rounded-xl">
              <button
                type="button"
                onClick={() => setDestination('queue')}
                className={`flex items-center justify-center gap-2 py-2 rounded-lg text-xs font-semibold transition ${
                  destination === 'queue'
                    ? 'bg-white text-sand-900 shadow-xs'
                    : 'text-sand-600 hover:text-sand-900'
                }`}
              >
                <Inbox className="w-3.5 h-3.5" />
                <span>Want to Consume (Queue)</span>
              </button>
              <button
                type="button"
                onClick={() => setDestination('shelf')}
                className={`flex items-center justify-center gap-2 py-2 rounded-lg text-xs font-semibold transition ${
                  destination === 'shelf'
                    ? 'bg-white text-sand-900 shadow-xs'
                    : 'text-sand-600 hover:text-sand-900'
                }`}
              >
                <Bookmark className="w-3.5 h-3.5" />
                <span>Already Finished (Shelf)</span>
              </button>
            </div>
          </div>

          {/* Media Category Pills */}
          <div>
            <label className="block text-xs font-semibold text-sand-700 uppercase tracking-wider mb-1.5">
              Category
            </label>
            <div className="grid grid-cols-5 gap-1.5">
              {MEDIA_TYPES.map((m) => (
                <button
                  type="button"
                  key={m.type}
                  onClick={() => setType(m.type)}
                  className={`py-2 px-1 rounded-xl text-center transition flex flex-col items-center gap-0.5 border ${
                    type === m.type
                      ? 'border-sand-800 bg-sand-900 text-white font-semibold shadow-xs'
                      : 'border-sand-200 bg-white text-sand-700 hover:bg-sand-50'
                  }`}
                >
                  <span className="text-sm">{m.icon}</span>
                  <span className="text-[10px] truncate max-w-full">{m.label}</span>
                </button>
              ))}
            </div>
          </div>

          {/* External URL */}
          <div>
            <div className="flex items-center justify-between mb-1">
              <label className="text-xs font-semibold text-sand-700 uppercase tracking-wider">
                External URL *
              </label>
              {isEnriching && (
                <span className="flex items-center gap-1 text-[11px] text-amber-700 font-medium">
                  <Loader2 className="w-3 h-3 animate-spin" />
                  <span>Extracting metadata...</span>
                </span>
              )}
              {enrichSuccess && !isEnriching && (
                <span className="flex items-center gap-1 text-[11px] text-emerald-700 font-medium">
                  <Sparkles className="w-3 h-3" />
                  <span>Auto-filled from source</span>
                </span>
              )}
            </div>
            <input
              type="url"
              required
              placeholder="e.g. Paste YouTube, TMDB, Open Library, or AniList URL"
              value={externalUrl}
              onChange={(e) => handleUrlChange(e.target.value)}
              className="w-full px-3.5 py-2 rounded-xl border border-sand-300 focus:border-sand-800 text-sm text-sand-900 outline-none transition bg-white"
            />
          </div>

          {/* Cover Art Live Preview if present */}
          {coverUrl && (
            <div className="p-3 bg-sand-50 rounded-xl border border-sand-200/80 flex items-center gap-3">
              <img
                src={coverUrl}
                alt="Preview"
                className="w-12 h-16 object-cover rounded-md border border-sand-200 flex-shrink-0 bg-sand-200"
                onError={(e) => {
                  (e.target as HTMLElement).style.display = 'none';
                }}
              />
              <div className="min-w-0 flex-1">
                <span className="text-[10px] uppercase font-bold text-sand-500 tracking-wider">
                  Metadata Preview
                </span>
                <p className="text-xs font-semibold text-sand-900 truncate">
                  {title || 'Loading title...'}
                </p>
                {creator && (
                  <p className="text-[11px] text-sand-600 truncate">
                    by {creator}
                  </p>
                )}
              </div>
            </div>
          )}

          {/* Title & Creator */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-semibold text-sand-700 uppercase tracking-wider mb-1">
                Title *
              </label>
              <input
                type="text"
                required
                placeholder="Title or work name"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full px-3.5 py-2 rounded-xl border border-sand-300 focus:border-sand-800 text-sm text-sand-900 outline-none transition bg-white"
              />
            </div>

            <div>
              <label className="block text-xs font-semibold text-sand-700 uppercase tracking-wider mb-1">
                Creator / Channel / Author
              </label>
              <input
                type="text"
                placeholder="e.g. Christopher Nolan, Denis Villeneuve"
                value={creator}
                onChange={(e) => setCreator(e.target.value)}
                className="w-full px-3.5 py-2 rounded-xl border border-sand-300 focus:border-sand-800 text-sm text-sand-900 outline-none transition bg-white"
              />
            </div>
          </div>

          {/* Cover Art Image URL */}
          <div>
            <label className="block text-xs font-semibold text-sand-700 uppercase tracking-wider mb-1">
              Cover Art URL (Optional)
            </label>
            <input
              type="url"
              placeholder="https://... cover image URL"
              value={coverUrl}
              onChange={(e) => setCoverUrl(e.target.value)}
              className="w-full px-3.5 py-2 rounded-xl border border-sand-300 focus:border-sand-800 text-sm text-sand-900 outline-none transition bg-white"
            />
          </div>

          {/* If saving directly to Shelf: Rating and Take */}
          {destination === 'shelf' && (
            <div className="pt-2 border-t border-sand-100 space-y-3">
              <div>
                <label className="block text-xs font-semibold text-sand-700 uppercase tracking-wider mb-1.5">
                  Rating ({rating}/10)
                </label>
                <div className="grid grid-cols-10 gap-1">
                  {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((num) => (
                    <button
                      type="button"
                      key={num}
                      onClick={() => setRating(num)}
                      className={`py-1.5 rounded-lg text-xs font-bold transition ${
                        rating === num
                          ? 'bg-amber-600 text-white shadow-xs'
                          : 'bg-sand-100 text-sand-600 hover:bg-sand-200'
                      }`}
                    >
                      {num}
                    </button>
                  ))}
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-sand-700 uppercase tracking-wider mb-1">
                  One-Line Take
                </label>
                <input
                  type="text"
                  maxLength={280}
                  placeholder="Why did you enjoy this? A quick personal thought."
                  value={take}
                  onChange={(e) => setTake(e.target.value)}
                  className="w-full px-3.5 py-2 rounded-xl border border-sand-300 focus:border-sand-800 text-sm text-sand-900 font-serif outline-none transition bg-white"
                />
              </div>
            </div>
          )}

          {error && (
            <p className="text-xs text-rose-600 bg-rose-50 border border-rose-200 p-2.5 rounded-lg">
              {error}
            </p>
          )}

          {/* Actions */}
          <div className="flex items-center justify-end gap-2.5 pt-2">
            <button
              type="button"
              onClick={onClose}
              disabled={isSubmitting}
              className="px-4 py-2 rounded-xl text-xs font-medium text-sand-600 hover:bg-sand-100 transition"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex items-center gap-1.5 px-5 py-2 bg-sand-900 hover:bg-sand-800 text-white rounded-xl text-xs font-semibold shadow-xs transition disabled:opacity-50"
            >
              <Plus className="w-3.5 h-3.5" />
              <span>{isSubmitting ? 'Saving...' : destination === 'queue' ? 'Add to Queue' : 'Save to Shelf'}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
