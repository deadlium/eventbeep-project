import React from "react";
import { Job } from "@/types/job";
import { ExternalLinkIcon, MapPinIcon, BuildingIcon, CalendarIcon } from "./Icons";

interface JobTableProps {
  jobs: Job[];
}

function getRoleBadgeStyle(role: string | null): string {
  if (!role) {
    return "bg-zinc-100 text-zinc-700 dark:bg-zinc-800 dark:text-zinc-300 border-zinc-200 dark:border-zinc-700";
  }
  const r = role.toLowerCase();
  if (r.includes("backend")) {
    return "bg-indigo-50 text-indigo-700 dark:bg-indigo-950/60 dark:text-indigo-300 border-indigo-200 dark:border-indigo-800/60";
  }
  if (r.includes("frontend")) {
    return "bg-sky-50 text-sky-700 dark:bg-sky-950/60 dark:text-sky-300 border-sky-200 dark:border-sky-800/60";
  }
  if (r.includes("full stack") || r.includes("fullstack")) {
    return "bg-purple-50 text-purple-700 dark:bg-purple-950/60 dark:text-purple-300 border-purple-200 dark:border-purple-800/60";
  }
  if (r.includes("devops") || r.includes("cloud") || r.includes("sre")) {
    return "bg-amber-50 text-amber-700 dark:bg-amber-950/60 dark:text-amber-300 border-amber-200 dark:border-amber-800/60";
  }
  if (r.includes("data") || r.includes("ai") || r.includes("ml")) {
    return "bg-emerald-50 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-300 border-emerald-200 dark:border-emerald-800/60";
  }
  if (r.includes("mobile") || r.includes("flutter") || r.includes("android") || r.includes("ios")) {
    return "bg-rose-50 text-rose-700 dark:bg-rose-950/60 dark:text-rose-300 border-rose-200 dark:border-rose-800/60";
  }
  if (r.includes("qa") || r.includes("test")) {
    return "bg-teal-50 text-teal-700 dark:bg-teal-950/60 dark:text-teal-300 border-teal-200 dark:border-teal-800/60";
  }
  return "bg-zinc-100 text-zinc-700 dark:bg-zinc-800 dark:text-zinc-300 border-zinc-200 dark:border-zinc-700";
}

function formatDate(isoDate: string): string {
  try {
    const d = new Date(isoDate);
    if (isNaN(d.getTime())) return "Recently";
    const now = new Date();
    const diffMs = now.getTime() - d.getTime();
    const diffHrs = Math.floor(diffMs / (1000 * 60 * 60));

    if (diffHrs < 1) return "Just now";
    if (diffHrs < 24) return `${diffHrs}h ago`;
    const diffDays = Math.floor(diffHrs / 24);
    if (diffDays < 7) return `${diffDays}d ago`;

    return d.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  } catch {
    return "Recently";
  }
}

export function JobTable({ jobs }: JobTableProps) {
  return (
    <div className="overflow-hidden rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
      {/* Desktop & Tablet Table */}
      <div className="hidden sm:block overflow-x-auto">
        <table className="min-w-full divide-y divide-zinc-200 dark:divide-zinc-800 text-left text-sm">
          <thead className="bg-zinc-50/70 dark:bg-zinc-800/40">
            <tr>
              <th
                scope="col"
                className="py-3.5 pl-6 pr-3 text-xs font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400"
              >
                Opportunity / Company
              </th>
              <th
                scope="col"
                className="px-4 py-3.5 text-xs font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400"
              >
                Role
              </th>
              <th
                scope="col"
                className="px-4 py-3.5 text-xs font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400"
              >
                Location
              </th>
              <th
                scope="col"
                className="px-4 py-3.5 text-xs font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400"
              >
                Discovered
              </th>
              <th
                scope="col"
                className="py-3.5 pl-3 pr-6 text-right text-xs font-semibold uppercase tracking-wider text-zinc-500 dark:text-zinc-400"
              >
                Action
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-200 dark:divide-zinc-800/70 bg-white dark:bg-zinc-900">
            {jobs.map((job) => (
              <tr
                key={job.id}
                className="group hover:bg-zinc-50/70 dark:hover:bg-zinc-800/30 transition-colors"
              >
                <td className="py-4 pl-6 pr-3 max-w-sm">
                  <div className="font-semibold text-zinc-900 dark:text-zinc-100 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors flex items-center gap-1.5">
                    <span className="truncate">{job.title}</span>
                  </div>
                  <div className="flex items-center gap-1 text-xs text-zinc-500 dark:text-zinc-400 mt-1">
                    <BuildingIcon className="w-3.5 h-3.5 text-zinc-400" />
                    <span className="font-medium text-zinc-700 dark:text-zinc-300">
                      {job.company}
                    </span>
                  </div>
                </td>
                <td className="px-4 py-4 whitespace-nowrap">
                  <span
                    className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium border ${getRoleBadgeStyle(
                      job.role
                    )}`}
                  >
                    {job.role || "General"}
                  </span>
                </td>
                <td className="px-4 py-4 whitespace-nowrap">
                  <span className="inline-flex items-center gap-1 text-xs text-zinc-600 dark:text-zinc-400">
                    <MapPinIcon className="w-3.5 h-3.5 text-zinc-400" />
                    {job.city || "Remote / Pan-India"}
                  </span>
                </td>
                <td className="px-4 py-4 whitespace-nowrap">
                  <span className="inline-flex items-center gap-1 text-xs text-zinc-500 dark:text-zinc-400">
                    <CalendarIcon className="w-3.5 h-3.5 text-zinc-400" />
                    {formatDate(job.created_at)}
                  </span>
                </td>
                <td className="py-4 pl-3 pr-6 text-right whitespace-nowrap">
                  <a
                    href={job.source_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1.5 rounded-lg border border-zinc-300 dark:border-zinc-700 bg-white dark:bg-zinc-800 px-3 py-1.5 text-xs font-semibold text-zinc-800 dark:text-zinc-200 shadow-sm hover:bg-zinc-50 dark:hover:bg-zinc-700 hover:text-indigo-600 dark:hover:text-indigo-400 hover:border-indigo-300 dark:hover:border-indigo-600 transition-all"
                  >
                    Apply
                    <ExternalLinkIcon className="h-3 w-3" />
                  </a>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Mobile Card List */}
      <div className="sm:hidden divide-y divide-zinc-200 dark:divide-zinc-800">
        {jobs.map((job) => (
          <div key={job.id} className="p-4 space-y-3">
            <div className="flex items-start justify-between gap-2">
              <div>
                <h4 className="text-sm font-semibold text-zinc-900 dark:text-zinc-100 leading-snug">
                  {job.title}
                </h4>
                <div className="flex items-center gap-1 text-xs text-zinc-500 dark:text-zinc-400 mt-0.5">
                  <BuildingIcon className="w-3.5 h-3.5" />
                  <span className="font-medium text-zinc-700 dark:text-zinc-300">
                    {job.company}
                  </span>
                </div>
              </div>
              <a
                href={job.source_url}
                target="_blank"
                rel="noopener noreferrer"
                className="shrink-0 inline-flex items-center gap-1 rounded-lg bg-indigo-600 px-2.5 py-1 text-xs font-semibold text-white shadow-sm hover:bg-indigo-500"
              >
                Apply
                <ExternalLinkIcon className="h-3 w-3" />
              </a>
            </div>

            <div className="flex flex-wrap items-center gap-2 text-xs">
              <span
                className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border ${getRoleBadgeStyle(
                  job.role
                )}`}
              >
                {job.role || "General"}
              </span>
              <span className="inline-flex items-center gap-1 text-zinc-500 dark:text-zinc-400">
                <MapPinIcon className="w-3.5 h-3.5" />
                {job.city || "Remote"}
              </span>
              <span className="text-zinc-300 dark:text-zinc-700">•</span>
              <span className="inline-flex items-center gap-1 text-zinc-400">
                <CalendarIcon className="w-3.5 h-3.5" />
                {formatDate(job.created_at)}
              </span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
