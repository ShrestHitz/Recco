import React, { useState } from 'react';
import { Check, X } from 'lucide-react';
import { ShelfEntry } from '../types';

interface FinishModalProps {
  entry: ShelfEntry | null;
  onClose: () => void;
  onSubmit: (rating: number, take: string) => Promise<void>;
}

const RATING_DESCRIPTIONS: Record<number, string> = {
  10: 'Masterpiece — Changed how I think',
  9: 'Incredible — Highly recommend to anyone',
  8: 'Great — Very much enjoyed',
  7: 'Good — Solid work with great moments',
  6: 'Decent — Worth a look if interested',
  5: 'Average — Has its merits and flaws',
  4: 'Subpar — Fell short of expectations',
  3: 'Weak — Hard to recommend',
  2: 'Poor — Didn’t work for me',
  1: 'Disliked — Would avoid',
};

export const FinishModal: React.FC<FinishModalProps> = ({ entry, onClose, onSubmit }) => {
  const [rating, setRating] = useState<number>(8);
  const [take, setTake] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  if (!entry || !entry.item) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (rating < 1 || rating > 10) {
      setError('Please select a rating between 1 and 10');
      return;
    }
    if (take.length > 280) {
      setError('Reflection must be under 280 characters');
      return;
    }

    try {
      setIsSubmitting(true);
      setError(null);
      await onSubmit(rating, take);
      onClose();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to finish item');
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
              Record to Taste Journal
            </h3>
            <p className="text-xs text-sand-500">
              Rate and leave a short personal reflection.
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-lg text-sand-400 hover:text-sand-700 hover:bg-sand-100 transition"
          >
            <X className="w-4 h-4" />
          </button>
        </div>

        {/* Content Details */}
        <form onSubmit={handleSubmit} className="p-6 space-y-5">
          <div className="flex items-center gap-3 p-3 rounded-xl bg-sand-50 border border-sand-200/80">
            {entry.item.cover_url ? (
              <img
                src={entry.item.cover_url}
                alt={entry.item.title}
                className="w-12 h-16 object-cover rounded-md flex-shrink-0 bg-sand-200"
                onError={(e) => {
                  (e.target as HTMLElement).style.display = 'none';
                }}
              />
            ) : null}
            <div className="min-w-0">
              <h4 className="text-sm font-semibold text-sand-900 truncate">
                {entry.item.title}
              </h4>
              {entry.item.creator && (
                <p className="text-xs text-sand-600 truncate">
                  by {entry.item.creator}
                </p>
              )}
            </div>
          </div>

          {/* 1 - 10 Rating Selector */}
          <div>
            <label className="block text-xs font-semibold text-sand-700 uppercase tracking-wider mb-2">
              Your Rating ({rating}/10)
            </label>
            <div className="grid grid-cols-10 gap-1 sm:gap-1.5">
              {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((num) => (
                <button
                  type="button"
                  key={num}
                  onClick={() => setRating(num)}
                  className={`py-2 rounded-lg text-xs font-bold transition flex flex-col items-center justify-center ${
                    rating === num
                      ? 'bg-amber-600 text-white shadow-xs scale-105'
                      : rating > num
                      ? 'bg-amber-100/70 text-amber-900 hover:bg-amber-200'
                      : 'bg-sand-100 text-sand-600 hover:bg-sand-200'
                  }`}
                >
                  <span>{num}</span>
                </button>
              ))}
            </div>
            <p className="text-xs text-amber-800 font-medium mt-2 text-center">
              {RATING_DESCRIPTIONS[rating]}
            </p>
          </div>

          {/* One-Line Take / Reflection */}
          <div>
            <div className="flex items-center justify-between mb-1.5">
              <label className="text-xs font-semibold text-sand-700 uppercase tracking-wider">
                One-Line Take (Optional)
              </label>
              <span className={`text-[11px] ${take.length > 250 ? 'text-amber-700 font-semibold' : 'text-sand-400'}`}>
                {take.length}/280
              </span>
            </div>
            <textarea
              value={take}
              onChange={(e) => setTake(e.target.value)}
              maxLength={280}
              rows={3}
              placeholder="e.g. Loved the character dynamics and the final scene. Truly unforgettable."
              className="w-full px-3.5 py-2.5 rounded-xl border border-sand-300 focus:border-sand-700 focus:ring-1 focus:ring-sand-700 text-sm text-sand-900 font-serif placeholder:font-sans placeholder:text-sand-400 outline-none transition resize-none bg-white"
            />
          </div>

          {error && (
            <p className="text-xs text-rose-600 bg-rose-50 border border-rose-200 p-2.5 rounded-lg">
              {error}
            </p>
          )}

          {/* Action buttons */}
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
              <Check className="w-3.5 h-3.5" />
              <span>{isSubmitting ? 'Recording...' : 'Add to Taste Journal'}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
