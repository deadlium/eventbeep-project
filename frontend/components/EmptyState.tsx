import React from "react";
import { SearchIcon } from "./Icons";

interface EmptyStateProps {
  onReset: () => void;
  hasFilters: boolean;
}

export function EmptyState({ onReset, hasFilters }: EmptyStateProps) {
  return (
    <div className="rounded-2xl border border-dashed border-zinc-300 dark:border-zinc-800 bg-white dark:bg-zinc-900/50 p-12 text-center shadow-sm">
      <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-zinc-100 dark:bg-zinc-800 text-zinc-400">
        <SearchIcon className="h-7 w-7" />
      </div>
      <h3 className="mt-4 text-base font-semibold text-zinc-900 dark:text-zinc-100">
        No jobs found matching your criteria
      </h3>
      <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400 max-w-md mx-auto">
        {hasFilters
          ? "Try adjusting your search terms, removing role or city filters, or searching with broader keywords."
          : "The pipeline has not aggregated any new jobs yet. The Go scraper runs periodically every 24 hours."}
      </p>
      {hasFilters && (
        <div className="mt-6">
          <button
            type="button"
            onClick={onReset}
            className="inline-flex items-center gap-2 rounded-xl bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 transition-colors"
          >
            Clear all filters
          </button>
        </div>
      )}
    </div>
  );
}
