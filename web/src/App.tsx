import React, { useEffect, useState, useCallback } from 'react';
import { Header } from './components/Header';
import { QueueView } from './components/QueueView';
import { ShelfView } from './components/ShelfView';
import { AddItemModal } from './components/AddItemModal';
import { FinishModal } from './components/FinishModal';
import { api } from './api/client';
import { ItemType, ShelfEntry } from './types';
import { Sparkles } from 'lucide-react';

const USER_ID_KEY = 'recco_demo_user_id';

function getOrCreateUserId(): string {
  let id = localStorage.getItem(USER_ID_KEY);
  if (!id) {
    // Generate valid UUID v4
    id = '10000000-1000-4000-8000-100000000001';
    localStorage.setItem(USER_ID_KEY, id);
  }
  return id;
}

export const App: React.FC = () => {
  const [userId] = useState<string>(getOrCreateUserId);
  const [activeTab, setActiveTab] = useState<'queue' | 'shelf'>('queue');
  const [queueEntries, setQueueEntries] = useState<ShelfEntry[]>([]);
  const [shelfEntries, setShelfEntries] = useState<ShelfEntry[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);
  const [isBackendHealthy, setIsBackendHealthy] = useState<boolean | null>(null);

  // Modals
  const [isAddModalOpen, setIsAddModalOpen] = useState<boolean>(false);
  const [finishModalEntry, setFinishModalEntry] = useState<ShelfEntry | null>(null);

  // Notification Toast
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

  const showToast = (message: string, type: 'success' | 'error' = 'success') => {
    setToast({ message, type });
    setTimeout(() => {
      setToast(null);
    }, 4000);
  };

  const loadData = useCallback(async () => {
    try {
      setIsRefreshing(true);

      // Check backend health
      const healthRes = await api.checkHealth();
      setIsBackendHealthy(healthRes.success && healthRes.data?.status === 'healthy');

      // Fetch queue and shelf concurrently
      const [queueRes, shelfRes] = await Promise.all([
        api.fetchQueue(userId),
        api.fetchShelf(userId),
      ]);

      if (queueRes.success && queueRes.data) {
        setQueueEntries(queueRes.data);
      }
      if (shelfRes.success && shelfRes.data) {
        setShelfEntries(shelfRes.data);
      }
    } catch (err: unknown) {
      console.error('Error fetching data:', err);
      setIsBackendHealthy(false);
    } finally {
      setLoading(false);
      setIsRefreshing(false);
    }
  }, [userId]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  // Handle Add Item (from AddItemModal)
  const handleAddItem = async (input: {
    type: ItemType;
    title: string;
    creator: string;
    external_url: string;
    cover_url?: string;
    destination: 'queue' | 'shelf';
    rating?: number;
    take?: string;
  }) => {
    // 1. Create item in catalog
    const itemRes = await api.createItem({
      type: input.type,
      title: input.title,
      creator: input.creator,
      external_url: input.external_url,
      cover_url: input.cover_url,
    });

    if (!itemRes.success || !itemRes.data) {
      throw new Error(itemRes.error || 'Failed to create item');
    }

    const createdItem = itemRes.data;

    // 2. Add to Queue or Shelf
    if (input.destination === 'queue') {
      const queueRes = await api.addToQueue(userId, createdItem.id);
      if (!queueRes.success) {
        throw new Error(queueRes.error || 'Failed to add to queue');
      }
      showToast(`"${createdItem.title}" added to your personal queue`);
      setActiveTab('queue');
    } else {
      const shelfRes = await api.finishItem(
        userId,
        createdItem.id,
        input.rating || 8,
        input.take || ''
      );
      if (!shelfRes.success) {
        throw new Error(shelfRes.error || 'Failed to record to shelf');
      }
      showToast(`"${createdItem.title}" saved to your taste journal`);
      setActiveTab('shelf');
    }

    await loadData();
  };

  // Handle Finish Item (from FinishModal)
  const handleFinishSubmit = async (rating: number, take: string) => {
    if (!finishModalEntry || !finishModalEntry.item) return;

    const res = await api.finishItem(userId, finishModalEntry.item_id, rating, take);
    if (!res.success) {
      throw new Error(res.error || 'Failed to finish item');
    }

    showToast(`Recorded "${finishModalEntry.item.title}" (${rating}/10) to your Taste Journal`);
    await loadData();
  };

  // Handle Drop Item
  const handleDrop = async (entry: ShelfEntry) => {
    if (!entry.item) return;
    const res = await api.dropItem(userId, entry.item_id);
    if (res.success) {
      showToast(`Marked "${entry.item.title}" as dropped`);
      await loadData();
    } else {
      showToast(res.error || 'Failed to drop item', 'error');
    }
  };

  // Handle Delete Entry
  const handleDelete = async (entry: ShelfEntry) => {
    const res = await api.deleteEntry(entry.id);
    if (res.success) {
      showToast('Entry removed');
      await loadData();
    } else {
      showToast(res.error || 'Failed to remove entry', 'error');
    }
  };

  // Handle Re-queue from Shelf
  const handleRequeue = async (entry: ShelfEntry) => {
    if (!entry.item) return;
    const res = await api.addToQueue(userId, entry.item_id);
    if (res.success) {
      showToast(`Returned "${entry.item.title}" to Queue`);
      setActiveTab('queue');
      await loadData();
    } else {
      showToast(res.error || 'Failed to return to queue', 'error');
    }
  };

  // Seed sample recommendations to test anti-distraction workflow
  const handleSeedStarterPack = async () => {
    try {
      setLoading(true);
      const sampleItems = [
        {
          type: 'movie' as ItemType,
          title: 'Interstellar',
          creator: 'Christopher Nolan',
          external_url: 'https://www.themoviedb.org/movie/157336',
          cover_url: 'https://image.tmdb.org/t/p/w500/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg',
        },
        {
          type: 'book' as ItemType,
          title: 'Designing Data-Intensive Applications',
          creator: 'Martin Kleppmann',
          external_url: 'https://openlibrary.org/works/OL17352697W',
          cover_url: 'https://covers.openlibrary.org/b/id/8091016-M.jpg',
        },
        {
          type: 'youtube' as ItemType,
          title: 'The Simplest Math Problem No One Can Solve',
          creator: 'Veritasium',
          external_url: 'https://www.youtube.com/watch?v=094y1Z2wpJg',
          cover_url: 'https://i.ytimg.com/vi/094y1Z2wpJg/hqdefault.jpg',
        },
      ];

      for (const sample of sampleItems) {
        const itemRes = await api.createItem(sample);
        if (itemRes.success && itemRes.data) {
          await api.addToQueue(userId, itemRes.data.id);
        }
      }

      showToast('Loaded 3 thoughtful starter recommendations into your queue');
      await loadData();
    } catch (err: unknown) {
      console.error('Failed seeding starter pack:', err);
      showToast('Could not seed starter pack. Please check backend connection.', 'error');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-sand-50 text-sand-900">
      {/* Header */}
      <Header
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        queueCount={queueEntries.length}
        shelfCount={shelfEntries.length}
        isBackendHealthy={isBackendHealthy}
        onOpenAddModal={() => setIsAddModalOpen(true)}
        onRefresh={loadData}
        isRefreshing={isRefreshing}
      />

      {/* Main Content Area */}
      <main className="flex-1 max-w-5xl w-full mx-auto px-4 sm:px-6 py-6 sm:py-8">
        
        {/* Offline Warning Banner */}
        {isBackendHealthy === false && (
          <div className="mb-6 p-4 rounded-2xl bg-amber-50 border border-amber-200 text-amber-900 text-xs sm:text-sm flex flex-col sm:flex-row sm:items-center justify-between gap-3 shadow-xs">
            <div>
              <p className="font-semibold">Go Backend API is offline</p>
              <p className="text-amber-800/90 text-xs mt-0.5">
                Run <code className="px-1.5 py-0.5 bg-amber-100 rounded text-amber-950 font-mono">go run ./cmd/api/main.go</code> in your terminal to connect the live API.
              </p>
            </div>
            <button
              onClick={loadData}
              className="px-3.5 py-1.5 bg-amber-200/80 hover:bg-amber-300 text-amber-950 font-medium rounded-xl text-xs transition self-start sm:self-auto"
            >
              Retry Connection
            </button>
          </div>
        )}

        {/* First-time starter pack banner (if both queue and shelf are completely empty) */}
        {!loading && queueEntries.length === 0 && shelfEntries.length === 0 && (
          <div className="mb-8 p-5 rounded-2xl bg-white border border-sand-200/90 shadow-xs flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="p-2.5 bg-amber-100 rounded-xl text-amber-900">
                <Sparkles className="w-5 h-5 text-amber-800" />
              </div>
              <div>
                <h4 className="text-sm font-semibold text-sand-900 font-serif">
                  Welcome to Recco — A quiet, finite culture space
                </h4>
                <p className="text-xs text-sand-600 mt-0.5">
                  Would you like to load a starter pack of 3 curated items to test the queue and shelf?
                </p>
              </div>
            </div>
            <button
              onClick={handleSeedStarterPack}
              className="px-4 py-2 bg-sand-900 hover:bg-sand-800 text-white rounded-xl text-xs font-semibold transition shadow-xs whitespace-nowrap"
            >
              Load Starter Pack
            </button>
          </div>
        )}

        {/* Active Tab View */}
        {activeTab === 'queue' ? (
          <QueueView
            entries={queueEntries}
            loading={loading}
            onFinishClick={(entry) => setFinishModalEntry(entry)}
            onDropClick={handleDrop}
            onDeleteClick={handleDelete}
            onOpenAddModal={() => setIsAddModalOpen(true)}
            onSwitchToShelf={() => setActiveTab('shelf')}
          />
        ) : (
          <ShelfView
            entries={shelfEntries}
            loading={loading}
            onRequeueClick={handleRequeue}
            onDeleteClick={handleDelete}
            onOpenAddModal={() => setIsAddModalOpen(true)}
            onSwitchToQueue={() => setActiveTab('queue')}
          />
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-sand-200/80 py-6 text-center text-xs text-sand-500">
        <p className="font-serif italic text-sand-600">
          "The best recommendations travel by word of mouth."
        </p>
        <p className="mt-1 text-[11px] text-sand-400">
          Recco Phase 1 — Anti-distraction, finite queue & personal taste journal.
        </p>
      </footer>

      {/* Toast Notification */}
      {toast && (
        <div className="fixed bottom-5 right-5 z-50 animate-in fade-in slide-in-from-bottom-5 duration-200">
          <div
            className={`px-4 py-3 rounded-2xl shadow-lg border text-xs sm:text-sm font-medium flex items-center gap-2 ${
              toast.type === 'success'
                ? 'bg-sand-900 text-white border-sand-800'
                : 'bg-rose-900 text-white border-rose-800'
            }`}
          >
            <span>{toast.message}</span>
          </div>
        </div>
      )}

      {/* Modals */}
      <AddItemModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        onSubmit={handleAddItem}
      />

      <FinishModal
        entry={finishModalEntry}
        onClose={() => setFinishModalEntry(null)}
        onSubmit={handleFinishSubmit}
      />
    </div>
  );
};
