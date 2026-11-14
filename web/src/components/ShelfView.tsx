import React, { useState } from 'react';
import { ExternalLink, Star, Trash2, Undo2 } from 'lucide-react';
import { ItemType, ShelfEntry } from '../types';
import { MediaBadge } from './MediaBadge';

interface ShelfViewProps {
  entries: ShelfEntry[];
  loading: boolean;
  onRequeueClick: (entry: ShelfEntry) => void;
  onDeleteClick: (entry: ShelfEntry) => void;
  onOpenAddModal: () => void;
  onSwitchToQueue: () => void;
}

const FILTER_TYPES: { label: string; value: ItemType | 'all' }[] = [
  { label: 'All Media', value: 'all' },
  { label: 'Movies & TV', value: 'movie' },
  { label: 'YouTube', value: 'youtube' },
  { label: 'Books', value: 'book' },
  { label: 'Anime', value: 'anime' },
  { label: 'Podcasts', value: 'podcast' },
];

export const ShelfView: React.FC<ShelfViewProps> = ({
  entries,
  loading,
  onRequeueClick,
  onDeleteClick,
  onOpenAddModal,
  onSwitchToQueue,
}) => {
  const [selectedFilter, setSelectedFilter] = useState<ItemType | 'all'>('all');

  const filteredEntries = entries.filter((e) => {
    if (selectedFilter === 'all') return true;
    return e.item?.type === selectedFilter;
  });

  if (loading) {
    return (
      <div className="py-20 flex flex-col items-center justify-center text-sand-500">
        <div className="w-8 h-8 border-2 border-sand-300 border-t-sand-800 rounded-full animate-spin mb-3" />
        <p className="text-sm">Loading your taste journal...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header and Filter Pills */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-3 border-b border-sand-200/60">
        <div>
          <h2 className="text-lg font-serif font-semibold text-sand-900">
            Taste Journal ({entries.length})
          </h2>
          <p className="text-xs text-sand-500">
            A permanent record of what moved, inspired, or informed you.
          </p>
        </div>

        {/* Filter Pills */}
        <div className="flex items-center gap-1.5 overflow-x-auto pb-1 sm:pb-0 scrollbar-none">
          {FILTER_TYPES.map((filter) => (
            <button
              key={filter.value}
              onClick={() => setSelectedFilter(filter.value)}
              className={`px-3 py-1 rounded-full text-xs font-medium whitespace-nowrap transition ${
                selectedFilter === filter.value
                  ? 'bg-sand-900 text-white font-semibold'
                  : 'bg-sand-100 text-sand-600 hover:bg-sand-200'
              }`}
            >
              {filter.label}
            </button>
          ))}
        </div>
      </div>

      {/* Empty State */}
      {filteredEntries.length === 0 ? (
        <div className="py-16 text-center max-w-sm mx-auto px-4">
          <div className="w-12 h-12 rounded-xl bg-sand-200/80 text-sand-600 flex items-center justify-center mx-auto mb-4">
            <Star className="w-6 h-6 text-sand-500" />
          </div>
          <h3 className="text-base font-serif font-semibold text-sand-900 mb-1">
            {selectedFilter === 'all'
              ? 'Your taste journal is empty'
              : `No finished ${selectedFilter} entries yet`}
          </h3>
          <p className="text-xs text-sand-500 mb-5">
            {selectedFilter === 'all'
              ? 'Finish items from your queue or directly record something you recently loved.'
              : 'Try selecting a different filter or finishing items in your queue.'}
          </p>
          <div className="flex items-center justify-center gap-2">
            <button
              onClick={onSwitchToQueue}
              className="px-4 py-2 bg-sand-900 text-white rounded-xl text-xs font-medium hover:bg-sand-800 transition"
            >
              Check Queue
            </button>
            <button
              onClick={onOpenAddModal}
              className="px-4 py-2 border border-sand-200 bg-white text-sand-700 rounded-xl text-xs font-medium hover:bg-sand-100 transition"
            >
              Record to Journal
            </button>
          </div>
        </div>
      ) : (
        /* Journal Grid */
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {filteredEntries.map((entry) => {
            const item = entry.item;
            if (!item) return null;

            return (
              <div
                key={entry.id}
                className="bg-white border border-sand-200/90 rounded-2xl p-4 sm:p-5 flex flex-col justify-between shadow-xs hover:border-sand-300 transition"
              >
                <div>
                  {/* Top Bar: Media Type, Rating, External Link */}
                  <div className="flex items-center justify-between gap-2 mb-3">
                    <MediaBadge type={item.type} />

                    <div className="flex items-center gap-2">
                      {entry.rating !== undefined && (
                        <div className="flex items-center gap-1 px-2.5 py-0.5 rounded-full bg-amber-50 border border-amber-200/70 text-amber-900 text-xs font-bold">
                          <Star className="w-3 h-3 fill-amber-500 text-amber-500" />
                          <span>{entry.rating}/10</span>
                        </div>
                      )}

                      <a
                        href={item.external_url}
                        target="_blank"
                        rel="noopener noreferrer"
                        title="Open external content"
                        className="p-1 text-sand-400 hover:text-sand-800 rounded transition"
                      >
                        <ExternalLink className="w-3.5 h-3.5" />
                      </a>
                    </div>
                  </div>

                  {/* Item Content Preview */}
                  <div className="flex gap-3.5">
                    {item.cover_url ? (
                      <img
                        src={item.cover_url}
                        alt={item.title}
                        className="w-14 h-20 object-cover rounded-lg border border-sand-200 flex-shrink-0 bg-sand-100"
                        onError={(e) => {
                          (e.target as HTMLElement).style.display = 'none';
                        }}
                      />
                    ) : (
                      <div className="w-14 h-20 rounded-lg border border-dashed border-sand-300 bg-sand-50 flex items-center justify-center flex-shrink-0 text-sand-400">
                        <MediaBadge type={item.type} showLabel={false} />
                      </div>
                    )}

                    <div className="flex-1 min-w-0">
                      <h3 className="text-base font-semibold text-sand-900 leading-snug line-clamp-2">
                        {item.title}
                      </h3>
                      {item.creator && (
                        <p className="text-xs text-sand-600 mt-0.5 line-clamp-1">
                          by <span className="font-medium">{item.creator}</span>
                        </p>
                      )}
                      {entry.finished_at && (
                        <p className="text-[11px] text-sand-400 mt-1">
                          Finished on {new Date(entry.finished_at).toLocaleDateString()}
                        </p>
                      )}
                    </div>
                  </div>

                  {/* One-Line Take (Personal Reflection) */}
                  {entry.take ? (
                    <div className="mt-3.5 p-3 rounded-xl bg-sand-50 border border-sand-200/70 text-xs sm:text-sm text-sand-800 font-serif italic leading-relaxed">
                      "{entry.take}"
                    </div>
                  ) : (
                    <div className="mt-3 text-[11px] text-sand-400 italic">
                      No personal note added.
                    </div>
                  )}
                </div>

                {/* Bottom Actions */}
                <div className="mt-4 pt-3 border-t border-sand-100 flex items-center justify-between">
                  <button
                    onClick={() => onRequeueClick(entry)}
                    className="inline-flex items-center gap-1 text-xs text-sand-500 hover:text-sand-800 transition"
                  >
                    <Undo2 className="w-3 h-3" />
                    <span>Return to Queue</span>
                  </button>

                  <button
                    onClick={() => onDeleteClick(entry)}
                    title="Remove from journal"
                    className="p-1 text-sand-400 hover:text-rose-600 rounded transition"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};
