import React from 'react';
import { BookOpen, Film, Headphones, PlaySquare, Sparkles } from 'lucide-react';
import { ItemType } from '../types';

interface MediaBadgeProps {
  type: ItemType;
  showLabel?: boolean;
}

export const MediaBadge: React.FC<MediaBadgeProps> = ({ type, showLabel = true }) => {
  switch (type) {
    case 'youtube':
      return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-medium bg-red-50 text-red-700 border border-red-200/60">
          <PlaySquare className="w-3 h-3" />
          {showLabel && <span>YouTube</span>}
        </span>
      );
    case 'movie':
      return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-medium bg-indigo-50 text-indigo-700 border border-indigo-200/60">
          <Film className="w-3 h-3" />
          {showLabel && <span>Movie / TV</span>}
        </span>
      );
    case 'book':
      return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-medium bg-amber-50 text-amber-800 border border-amber-200/60">
          <BookOpen className="w-3 h-3" />
          {showLabel && <span>Book</span>}
        </span>
      );
    case 'anime':
      return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-medium bg-purple-50 text-purple-700 border border-purple-200/60">
          <Sparkles className="w-3 h-3" />
          {showLabel && <span>Anime</span>}
        </span>
      );
    case 'podcast':
      return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-medium bg-emerald-50 text-emerald-700 border border-emerald-200/60">
          <Headphones className="w-3 h-3" />
          {showLabel && <span>Podcast</span>}
        </span>
      );
    default:
      return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[11px] font-medium bg-sand-100 text-sand-700 border border-sand-200">
          <span>{type}</span>
        </span>
      );
  }
};
