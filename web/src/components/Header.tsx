import React from 'react';
import { Bookmark, Inbox, Plus, RefreshCw } from 'lucide-react';

interface HeaderProps {
  activeTab: 'queue' | 'shelf';
  setActiveTab: (tab: 'queue' | 'shelf') => void;
  queueCount: number;
  shelfCount: number;
  isBackendHealthy: boolean | null;
  onOpenAddModal: () => void;
  onRefresh: () => void;
  isRefreshing: boolean;
}

export const Header: React.FC<HeaderProps> = ({
  activeTab,
  setActiveTab,
  queueCount,
  shelfCount,
  isBackendHealthy,
  onOpenAddModal,
  onRefresh,
  isRefreshing,
}) => {
  return (
    <header className="border-b border-sand-200/80 bg-sand-50/90 backdrop-blur sticky top-0 z-30">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 py-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          
          {/* Logo & Vision */}
          <div>
            <div className="flex items-center gap-2.5">
              <span className="text-2xl" role="img" aria-label="couch">🛋️</span>
              <h1 className="text-xl sm:text-2xl font-bold tracking-tight text-sand-900 font-serif">
                Recco
              </h1>
              <div className="flex items-center gap-1.5 ml-2 px-2 py-0.5 rounded-full text-xs font-medium border border-sand-200 bg-sand-100 text-sand-700">
                <span
                  className={`w-1.5 h-1.5 rounded-full ${
                    isBackendHealthy === true
                      ? 'bg-emerald-500 animate-pulse'
                      : isBackendHealthy === false
                      ? 'bg-rose-500'
                      : 'bg-amber-400'
                  }`}
                />
                <span>{isBackendHealthy ? 'Go API Online' : isBackendHealthy === false ? 'Offline' : 'Connecting...'}</span>
              </div>
            </div>
            <p className="text-xs sm:text-sm text-sand-600 mt-0.5">
              Recommended by people, not algorithms. <span className="hidden sm:inline">No feed. Just friends.</span>
            </p>
          </div>

          {/* Actions & Tab Switcher */}
          <div className="flex items-center gap-2 sm:gap-3 self-stretch sm:self-auto justify-between sm:justify-end">
            
            {/* Navigation Tabs */}
            <div className="flex p-1 bg-sand-200/70 rounded-xl text-xs sm:text-sm font-medium">
              <button
                onClick={() => setActiveTab('queue')}
                className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg transition-all ${
                  activeTab === 'queue'
                    ? 'bg-white text-sand-900 shadow-sm font-semibold'
                    : 'text-sand-700 hover:text-sand-900 hover:bg-sand-200/50'
                }`}
              >
                <Inbox className="w-3.5 h-3.5" />
                <span>Queue</span>
                <span className={`px-1.5 py-0.2 rounded-full text-[11px] ${
                  activeTab === 'queue' ? 'bg-amber-100 text-amber-900' : 'bg-sand-300/60 text-sand-700'
                }`}>
                  {queueCount}
                </span>
              </button>

              <button
                onClick={() => setActiveTab('shelf')}
                className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg transition-all ${
                  activeTab === 'shelf'
                    ? 'bg-white text-sand-900 shadow-sm font-semibold'
                    : 'text-sand-700 hover:text-sand-900 hover:bg-sand-200/50'
                }`}
              >
                <Bookmark className="w-3.5 h-3.5" />
                <span>Taste Journal</span>
                <span className={`px-1.5 py-0.2 rounded-full text-[11px] ${
                  activeTab === 'shelf' ? 'bg-amber-100 text-amber-900' : 'bg-sand-300/60 text-sand-700'
                }`}>
                  {shelfCount}
                </span>
              </button>
            </div>

            {/* Refresh */}
            <button
              onClick={onRefresh}
              disabled={isRefreshing}
              title="Refresh queue and shelf"
              className="p-2 rounded-xl border border-sand-200 bg-white hover:bg-sand-100 text-sand-700 transition"
            >
              <RefreshCw className={`w-4 h-4 ${isRefreshing ? 'animate-spin' : ''}`} />
            </button>

            {/* Add Item Button */}
            <button
              onClick={onOpenAddModal}
              className="flex items-center gap-1.5 px-3.5 py-2 bg-sand-900 hover:bg-sand-800 text-white rounded-xl text-xs sm:text-sm font-medium transition shadow-sm active:scale-95"
            >
              <Plus className="w-4 h-4" />
              <span>Save Item</span>
            </button>

          </div>

        </div>
      </div>
    </header>
  );
};
