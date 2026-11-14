import React from 'react';
import { Check, ExternalLink, Sparkles, Trash2, XCircle } from 'lucide-react';
import { ShelfEntry } from '../types';
import { MediaBadge } from './MediaBadge';

interface QueueViewProps {
  entries: ShelfEntry[];
  loading: boolean;
  onFinishClick: (entry: ShelfEntry) => void;
  onDropClick: (entry: ShelfEntry) => void;
  onDeleteClick: (entry: ShelfEntry) => void;
  onOpenAddModal: () => void;
  onSwitchToShelf: () => void;
}

export const QueueView: React.FC<QueueViewProps> = ({
  entries,
  loading,
  onFinishClick,
  onDropClick,
  onDeleteClick,
  onOpenAddModal,
  onSwitchToShelf,
}) => {
  if (loading) {
    return (
      <div className="py-20 flex flex-col items-center justify-center text-sand-500">
        <div className="w-8 h-8 border-2 border-sand-300 border-t-sand-800 rounded-full animate-spin mb-3" />
        <p className="text-sm">Retrieving your personal queue...</p>
      </div>
    );
  }

  // Queue Zero State (Anti-distraction Core Feature)
  if (entries.length === 0) {
    return (
      <div className="py-16 sm:py-24 text-center max-w-lg mx-auto px-4">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-amber-100/70 text-amber-900 mb-6 shadow-inner">
          <Sparkles className="w-8 h-8 text-amber-800" />
        </div>

        <h2 className="text-2xl sm:text-3xl font-serif font-bold text-sand-900 tracking-tight mb-3">
          You've reached Queue Zero.
        </h2>

        <p className="text-sand-600 text-sm sm:text-base leading-relaxed mb-8">
          No algorithmic feed pushing endless content. You’ve finished what was on your mind. Close the tab, go for a walk, or read something offline.
        </p>

        <div className="flex flex-col sm:flex-row items-center justify-center gap-3">
          <button
            onClick={onOpenAddModal}
            className="w-full sm:w-auto px-5 py-2.5 bg-sand-900 hover:bg-sand-800 text-white rounded-xl text-sm font-medium transition shadow-sm"
          >
            Queue New Recommendation
          </button>
          <button
            onClick={onSwitchToShelf}
            className="w-full sm:w-auto px-5 py-2.5 bg-white border border-sand-200 hover:bg-sand-100 text-sand-800 rounded-xl text-sm font-medium transition"
          >
            View Taste Journal
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between pb-2 border-b border-sand-200/60">
        <div>
          <h2 className="text-lg font-serif font-semibold text-sand-900">
            Personal Queue ({entries.length})
          </h2>
          <p className="text-xs text-sand-500">
            Curated things you want to watch, read, or listen to. Links open directly in source apps.
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {entries.map((entry) => {
          const item = entry.item;
          if (!item) return null;

          return (
            <div
              key={entry.id}
              className="bg-white border border-sand-200/90 rounded-2xl p-4 sm:p-5 flex flex-col justify-between shadow-xs hover:border-sand-300 transition duration-200 group"
            >
              <div>
                <div className="flex items-start justify-between gap-3 mb-2.5">
                  <MediaBadge type={item.type} />
                  <a
                    href={item.external_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    title="Open in native application (no embedded autoplay)"
                    className="inline-flex items-center gap-1 text-xs font-medium text-sand-500 hover:text-sand-900 transition px-2 py-0.5 rounded-md hover:bg-sand-100"
                  >
                    <span>Open</span>
                    <ExternalLink className="w-3 h-3" />
                  </a>
                </div>

                <div className="flex gap-4">
                  {item.cover_url ? (
                    <img
                      src={item.cover_url}
                      alt={item.title}
                      className="w-16 h-22 object-cover rounded-lg border border-sand-200 shadow-xs flex-shrink-0 bg-sand-100"
                      onError={(e) => {
                        (e.target as HTMLElement).style.display = 'none';
                      }}
                    />
                  ) : (
                    <div className="w-16 h-22 rounded-lg border border-dashed border-sand-300 bg-sand-50 flex items-center justify-center flex-shrink-0 text-sand-400">
                      <MediaBadge type={item.type} showLabel={false} />
                    </div>
                  )}

                  <div className="flex-1 min-w-0">
                    <h3 className="text-base font-semibold text-sand-900 leading-snug line-clamp-2">
                      {item.title}
                    </h3>
                    {item.creator && (
                      <p className="text-xs text-sand-600 mt-1 line-clamp-1">
                        by <span className="font-medium">{item.creator}</span>
                      </p>
                    )}
                    <p className="text-[11px] text-sand-400 mt-2">
                      Added {new Date(entry.created_at).toLocaleDateString()}
                    </p>
                  </div>
                </div>
              </div>

              {/* Action buttons */}
              <div className="mt-4 pt-3 border-t border-sand-100 flex items-center justify-between gap-2">
                <div className="flex items-center gap-1">
                  <button
                    onClick={() => onDropClick(entry)}
                    title="Drop item (discontinued)"
                    className="p-1.5 text-sand-400 hover:text-sand-700 hover:bg-sand-100 rounded-lg text-xs transition"
                  >
                    <XCircle className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => onDeleteClick(entry)}
                    title="Remove from queue"
                    className="p-1.5 text-sand-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg text-xs transition"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>

                <button
                  onClick={() => onFinishClick(entry)}
                  className="flex items-center gap-1.5 px-3 py-1.5 bg-amber-600 hover:bg-amber-700 text-white rounded-xl text-xs font-semibold shadow-xs transition active:scale-95"
                >
                  <Check className="w-3.5 h-3.5" />
                  <span>Finished & Rate</span>
                </button>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
