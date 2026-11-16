import React from "react";

export function JobTableSkeleton({ rows = 6 }: { rows?: number }) {
  return (
    <div className="overflow-hidden rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
      <div className="hidden sm:block overflow-x-auto">
        <table className="min-w-full divide-y divide-zinc-200 dark:divide-zinc-800 text-left text-sm">
          <thead className="bg-zinc-50/80 dark:bg-zinc-800/50">
            <tr>
              <th className="py-3.5 pl-6 pr-3 text-xs font-semibold text-zinc-500 uppercase tracking-wider">
                Job & Company
              </th>
              <th className="px-4 py-3.5 text-xs font-semibold text-zinc-500 uppercase tracking-wider">
                Role
              </th>
              <th className="px-4 py-3.5 text-xs font-semibold text-zinc-500 uppercase tracking-wider">
                Location
              </th>
              <th className="px-4 py-3.5 text-xs font-semibold text-zinc-500 uppercase tracking-wider">
                Discovered
              </th>
              <th className="py-3.5 pl-3 pr-6 text-right text-xs font-semibold text-zinc-500 uppercase tracking-wider">
                Action
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800 bg-white dark:bg-zinc-900">
            {Array.from({ length: rows }).map((_, i) => (
              <tr key={i} className="animate-pulse">
                <td className="py-4 pl-6 pr-3">
                  <div className="h-4 w-48 rounded bg-zinc-200 dark:bg-zinc-800 mb-2" />
                  <div className="h-3 w-28 rounded bg-zinc-100 dark:bg-zinc-800/60" />
                </td>
                <td className="px-4 py-4">
                  <div className="h-6 w-24 rounded-full bg-zinc-200 dark:bg-zinc-800" />
                </td>
                <td className="px-4 py-4">
                  <div className="h-4 w-20 rounded bg-zinc-200 dark:bg-zinc-800" />
                </td>
                <td className="px-4 py-4">
                  <div className="h-4 w-16 rounded bg-zinc-200 dark:bg-zinc-800" />
                </td>
                <td className="py-4 pl-3 pr-6 text-right">
                  <div className="inline-block h-8 w-20 rounded-lg bg-zinc-200 dark:bg-zinc-800" />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Mobile Card Skeleton */}
      <div className="sm:hidden divide-y divide-zinc-200 dark:divide-zinc-800 p-4 space-y-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="animate-pulse pt-3 space-y-3">
            <div className="h-4 w-3/4 rounded bg-zinc-200 dark:bg-zinc-800" />
            <div className="h-3 w-1/3 rounded bg-zinc-100 dark:bg-zinc-800/60" />
            <div className="flex gap-2">
              <div className="h-5 w-20 rounded-full bg-zinc-200 dark:bg-zinc-800" />
              <div className="h-5 w-16 rounded-full bg-zinc-200 dark:bg-zinc-800" />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
