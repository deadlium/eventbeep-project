"use client";

import React, { useState, useEffect } from "react";
import { Header } from "@/components/Header";
import { SearchFilters } from "@/components/SearchFilters";
import { JobTable } from "@/components/JobTable";
import { JobTableSkeleton } from "@/components/JobSkeleton";
import { EmptyState } from "@/components/EmptyState";
import { ErrorBanner } from "@/components/ErrorBanner";
import { PaginationControls } from "@/components/PaginationControls";
import { useDebounce } from "@/hooks/useDebounce";
import { fetchJobs } from "@/lib/api";
import { Job, Pagination } from "@/types/job";
import { SparklesIcon, BriefcaseIcon, BuildingIcon } from "@/components/Icons";

export default function HomePage() {
  const [searchTerm, setSearchTerm] = useState("");
  const debouncedSearch = useDebounce(searchTerm, 350);

  const [selectedRole, setSelectedRole] = useState("All Roles");
  const [selectedCity, setSelectedCity] = useState("All Cities");
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(20);

  const [jobs, setJobs] = useState<Job[]>([]);
  const [pagination, setPagination] = useState<Pagination>({
    total: 0,
    page: 1,
    limit: 20,
    total_pages: 1,
    has_next: false,
    has_prev: false,
  });

  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isFallback, setIsFallback] = useState(false);
  const [reloadTrigger, setReloadTrigger] = useState(0);

  // Trigger data loading on filter/page/reload change
  useEffect(() => {
    let ignore = false;

    async function loadData() {
      try {
        const { response, isFallback: fallbackUsed } = await fetchJobs({
          page,
          limit,
          role: selectedRole,
          city: selectedCity,
          search: debouncedSearch,
        });

        if (!ignore) {
          setJobs(response.data);
          setPagination(response.pagination);
          setIsFallback(fallbackUsed);
          setError(null);
          setIsLoading(false);
        }
      } catch (err: unknown) {
        if (!ignore) {
          const message =
            err instanceof Error ? err.message : "Failed to load jobs from API";
          setError(message);
          setIsLoading(false);
        }
      }
    }

    loadData();

    return () => {
      ignore = true;
    };
  }, [page, limit, selectedRole, selectedCity, debouncedSearch, reloadTrigger]);

  const handleSearchChange = (val: string) => {
    setIsLoading(true);
    setSearchTerm(val);
    setPage(1);
  };

  const handleRoleChange = (role: string) => {
    setIsLoading(true);
    setSelectedRole(role);
    setPage(1);
  };

  const handleCityChange = (city: string) => {
    setIsLoading(true);
    setSelectedCity(city);
    setPage(1);
  };

  const handleLimitChange = (newLimit: number) => {
    setIsLoading(true);
    setLimit(newLimit);
    setPage(1);
  };

  const handlePageChange = (newPage: number) => {
    setIsLoading(true);
    setPage(newPage);
  };

  const handleResetFilters = () => {
    setIsLoading(true);
    setSearchTerm("");
    setSelectedRole("All Roles");
    setSelectedCity("All Cities");
    setPage(1);
  };

  const handleRetry = () => {
    setIsLoading(true);
    setError(null);
    setReloadTrigger((prev) => prev + 1);
  };

  const hasActiveFilters =
    searchTerm.trim() !== "" ||
    selectedRole !== "All Roles" ||
    selectedCity !== "All Cities";

  return (
    <div className="min-h-screen flex flex-col bg-zinc-50 dark:bg-zinc-950 font-sans text-zinc-900 dark:text-zinc-100">
      <Header isFallback={isFallback} totalJobs={pagination.total} />

      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 space-y-6">
        {/* Hero Section */}
        <div className="relative overflow-hidden rounded-3xl bg-linear-to-r from-indigo-900 via-zinc-900 to-zinc-900 border border-indigo-950/60 p-6 sm:p-8 shadow-xl">
          <div className="relative z-10 max-w-3xl space-y-3">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-indigo-500/10 border border-indigo-400/20 text-indigo-300 text-xs font-semibold">
              <SparklesIcon className="w-3.5 h-3.5" />
              Automated Aggregation Engine • India Edition
            </div>
            <h2 className="text-2xl sm:text-3xl lg:text-4xl font-extrabold tracking-tight text-white">
              Discover Verified Fresher Jobs & Internships
            </h2>
            <p className="text-sm sm:text-base text-zinc-300 leading-relaxed">
              Scraped continuously from India&apos;s leading tech job portals, deduplicated with SHA-256 signatures, and enriched via Google Gemini LLM for precise role and location categorization.
            </p>
            <div className="flex flex-wrap items-center gap-4 pt-2 text-xs text-zinc-400">
              <span className="flex items-center gap-1.5">
                <BriefcaseIcon className="w-3.5 h-3.5 text-indigo-400" />
                Engineering, DevOps, AI & Product
              </span>
              <span>•</span>
              <span className="flex items-center gap-1.5">
                <BuildingIcon className="w-3.5 h-3.5 text-emerald-400" />
                Bangalore, Pune, Hyderabad, Remote
              </span>
            </div>
          </div>
          {/* Subtle decorative background blur */}
          <div className="absolute -right-16 -top-16 w-80 h-80 bg-indigo-600/20 rounded-full blur-3xl pointer-events-none" />
        </div>

        {/* Filter & Search Bar */}
        <SearchFilters
          searchTerm={searchTerm}
          onSearchChange={handleSearchChange}
          selectedRole={selectedRole}
          onRoleChange={handleRoleChange}
          selectedCity={selectedCity}
          onCityChange={handleCityChange}
          limit={limit}
          onLimitChange={handleLimitChange}
          onReset={handleResetFilters}
          hasActiveFilters={hasActiveFilters}
        />

        {/* Error Notification */}
        {error && (
          <ErrorBanner message={error} onRetry={handleRetry} />
        )}

        {/* Jobs Listing / Skeletons / Empty State */}
        {isLoading ? (
          <JobTableSkeleton rows={limit > 10 ? 10 : limit} />
        ) : jobs.length === 0 ? (
          <EmptyState
            hasFilters={hasActiveFilters}
            onReset={handleResetFilters}
          />
        ) : (
          <div className="space-y-4">
            <JobTable jobs={jobs} />
            <PaginationControls
              pagination={pagination}
              onPageChange={handlePageChange}
            />
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-950 py-6 mt-12 text-center text-xs text-zinc-500 dark:text-zinc-400">
        <div className="max-w-7xl mx-auto px-4 space-y-1">
          <p className="font-medium text-zinc-700 dark:text-zinc-300">
            EventBeep • Next.js & TypeScript UI • Phase 4 Architecture
          </p>
          <p>
            Powered by Go 1.23 REST API, PostgreSQL, and Gemini LLM Entity Enrichment.
          </p>
        </div>
      </footer>
    </div>
  );
}
