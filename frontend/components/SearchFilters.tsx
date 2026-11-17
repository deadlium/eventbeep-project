import React from "react";
import { SearchIcon, XIcon, BriefcaseIcon, MapPinIcon } from "./Icons";

export const ROLE_OPTIONS = [
  "All Roles",
  "Backend Engineer",
  "Frontend Developer",
  "Full Stack",
  "DevOps",
  "Data / AI",
  "Mobile",
  "QA / Testing",
];

export const CITY_OPTIONS = [
  "All Cities",
  "Bangalore",
  "Pune",
  "Hyderabad",
  "Mumbai",
  "Delhi NCR",
  "Chennai",
  "Remote",
];

interface SearchFiltersProps {
  searchTerm: string;
  onSearchChange: (value: string) => void;
  selectedRole: string;
  onRoleChange: (role: string) => void;
  selectedCity: string;
  onCityChange: (city: string) => void;
  limit: number;
  onLimitChange: (limit: number) => void;
  onReset: () => void;
  hasActiveFilters: boolean;
}

export function SearchFilters({
  searchTerm,
  onSearchChange,
  selectedRole,
  onRoleChange,
  selectedCity,
  onCityChange,
  limit,
  onLimitChange,
  onReset,
  hasActiveFilters,
}: SearchFiltersProps) {
  return (
    <div className="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4 sm:p-5 shadow-sm space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-12 gap-3 sm:gap-4">
        {/* Search Input */}
        <div className="md:col-span-5 relative">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5 text-zinc-400">
            <SearchIcon className="h-5 w-5" />
          </div>
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => onSearchChange(e.target.value)}
            placeholder="Search title, skills, company (e.g. 'frotend', 'pythn', 'banglore')..."
            className="w-full rounded-xl border border-zinc-300 dark:border-zinc-700 bg-zinc-50/50 dark:bg-zinc-800/60 pl-10 pr-9 py-2.5 text-sm text-zinc-900 dark:text-zinc-100 placeholder-zinc-400 focus:border-indigo-500 focus:bg-white dark:focus:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all"
          />
          {searchTerm && (
            <button
              onClick={() => onSearchChange("")}
              className="absolute inset-y-0 right-0 flex items-center pr-3 text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-200"
              title="Clear search"
            >
              <XIcon className="h-4 w-4" />
            </button>
          )}
        </div>

        {/* Role Selector */}
        <div className="md:col-span-3 relative">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-zinc-400">
            <BriefcaseIcon className="h-4 w-4" />
          </div>
          <select
            value={selectedRole}
            onChange={(e) => onRoleChange(e.target.value)}
            className="w-full appearance-none rounded-xl border border-zinc-300 dark:border-zinc-700 bg-zinc-50/50 dark:bg-zinc-800/60 pl-9 pr-8 py-2.5 text-sm text-zinc-900 dark:text-zinc-100 focus:border-indigo-500 focus:bg-white dark:focus:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer"
          >
            {ROLE_OPTIONS.map((role) => (
              <option key={role} value={role} className="dark:bg-zinc-900">
                {role === "All Roles" ? "All Engineering Roles" : role}
              </option>
            ))}
          </select>
          <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3 text-zinc-400">
            <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="m19 9-7 7-7-7" />
            </svg>
          </div>
        </div>

        {/* City Selector */}
        <div className="md:col-span-3 relative">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3 text-zinc-400">
            <MapPinIcon className="h-4 w-4" />
          </div>
          <select
            value={selectedCity}
            onChange={(e) => onCityChange(e.target.value)}
            className="w-full appearance-none rounded-xl border border-zinc-300 dark:border-zinc-700 bg-zinc-50/50 dark:bg-zinc-800/60 pl-9 pr-8 py-2.5 text-sm text-zinc-900 dark:text-zinc-100 focus:border-indigo-500 focus:bg-white dark:focus:bg-zinc-900 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer"
          >
            {CITY_OPTIONS.map((city) => (
              <option key={city} value={city} className="dark:bg-zinc-900">
                {city === "All Cities" ? "All Indian Tech Hubs" : city}
              </option>
            ))}
          </select>
          <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3 text-zinc-400">
            <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="m19 9-7 7-7-7" />
            </svg>
          </div>
        </div>

        {/* Per-Page Selector */}
        <div className="md:col-span-1">
          <select
            value={limit}
            onChange={(e) => onLimitChange(Number(e.target.value))}
            className="w-full rounded-xl border border-zinc-300 dark:border-zinc-700 bg-zinc-50/50 dark:bg-zinc-800/60 px-2 py-2.5 text-sm text-zinc-900 dark:text-zinc-100 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer text-center"
            title="Items per page"
          >
            <option value={10}>10</option>
            <option value={20}>20</option>
            <option value={50}>50</option>
          </select>
        </div>
      </div>

      {/* Active filters pill list & Reset button */}
      {hasActiveFilters && (
        <div className="flex flex-wrap items-center gap-2 pt-1 border-t border-zinc-100 dark:border-zinc-800/60 text-xs text-zinc-500">
          <span className="font-medium text-zinc-600 dark:text-zinc-400">Active filters:</span>
          {searchTerm && (
            <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-md bg-indigo-50 text-indigo-700 dark:bg-indigo-950/40 dark:text-indigo-300 border border-indigo-200 dark:border-indigo-800/60">
              Keyword: &quot;{searchTerm}&quot;
              <button onClick={() => onSearchChange("")} className="hover:text-indigo-900 ml-0.5">
                <XIcon className="w-3 h-3" />
              </button>
            </span>
          )}
          {selectedRole !== "All Roles" && (
            <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-md bg-violet-50 text-violet-700 dark:bg-violet-950/40 dark:text-violet-300 border border-violet-200 dark:border-violet-800/60">
              Role: {selectedRole}
              <button onClick={() => onRoleChange("All Roles")} className="hover:text-violet-900 ml-0.5">
                <XIcon className="w-3 h-3" />
              </button>
            </span>
          )}
          {selectedCity !== "All Cities" && (
            <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-md bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800/60">
              City: {selectedCity}
              <button onClick={() => onCityChange("All Cities")} className="hover:text-emerald-900 ml-0.5">
                <XIcon className="w-3 h-3" />
              </button>
            </span>
          )}
          <button
            onClick={onReset}
            className="ml-auto text-xs font-semibold text-rose-600 dark:text-rose-400 hover:underline hover:text-rose-700 flex items-center gap-1"
          >
            Reset All
          </button>
        </div>
      )}
    </div>
  );
}
