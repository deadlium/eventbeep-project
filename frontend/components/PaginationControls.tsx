import React from "react";
import { ChevronLeftIcon, ChevronRightIcon } from "./Icons";
import { Pagination } from "@/types/job";

interface PaginationControlsProps {
  pagination: Pagination;
  onPageChange: (newPage: number) => void;
}

export function PaginationControls({
  pagination,
  onPageChange,
}: PaginationControlsProps) {
  const { page, total_pages, total, limit, has_next, has_prev } = pagination;

  if (total === 0) return null;

  const start = Math.min((page - 1) * limit + 1, total);
  const end = Math.min(page * limit, total);

  // Helper to generate page buttons
  const getPageNumbers = () => {
    const pages: (number | string)[] = [];
    const maxVisible = 5;

    if (total_pages <= maxVisible) {
      for (let i = 1; i <= total_pages; i++) {
        pages.push(i);
      }
    } else {
      if (page <= 3) {
        pages.push(1, 2, 3, 4, "...", total_pages);
      } else if (page >= total_pages - 2) {
        pages.push(1, "...", total_pages - 3, total_pages - 2, total_pages - 1, total_pages);
      } else {
        pages.push(1, "...", page - 1, page, page + 1, "...", total_pages);
      }
    }
    return pages;
  };

  return (
    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 py-4 px-2">
      <div className="text-xs text-zinc-500 dark:text-zinc-400">
        Showing <span className="font-semibold text-zinc-900 dark:text-zinc-100">{start}</span> to{" "}
        <span className="font-semibold text-zinc-900 dark:text-zinc-100">{end}</span> of{" "}
        <span className="font-semibold text-zinc-900 dark:text-zinc-100">{total}</span> fresher opportunities
      </div>

      <div className="flex items-center gap-1.5">
        <button
          onClick={() => onPageChange(page - 1)}
          disabled={!has_prev}
          className="inline-flex items-center gap-1 px-3 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 text-xs font-medium text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
        >
          <ChevronLeftIcon className="w-3.5 h-3.5" />
          Previous
        </button>

        <div className="hidden sm:flex items-center gap-1">
          {getPageNumbers().map((p, idx) =>
            typeof p === "number" ? (
              <button
                key={idx}
                onClick={() => onPageChange(p)}
                className={`min-w-[32px] h-8 px-2 rounded-lg text-xs font-medium transition-all ${
                  p === page
                    ? "bg-indigo-600 text-white font-semibold shadow-sm"
                    : "border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800"
                }`}
              >
                {p}
              </button>
            ) : (
              <span key={idx} className="px-1 text-xs text-zinc-400">
                {p}
              </span>
            )
          )}
        </div>

        <button
          onClick={() => onPageChange(page + 1)}
          disabled={!has_next}
          className="inline-flex items-center gap-1 px-3 py-1.5 rounded-lg border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 text-xs font-medium text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800 disabled:opacity-40 disabled:cursor-not-allowed transition-all"
        >
          Next
          <ChevronRightIcon className="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  );
}
