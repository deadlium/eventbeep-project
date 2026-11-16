import React from "react";
import { SparklesIcon } from "./Icons";

interface HeaderProps {
  isFallback: boolean;
  totalJobs: number;
}

export function Header({ isFallback, totalJobs }: HeaderProps) {
  return (
    <header className="border-b border-zinc-200 bg-white/80 backdrop-blur-md sticky top-0 z-30 dark:border-zinc-800 dark:bg-zinc-950/80">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-xl bg-gradient-to-tr from-indigo-600 via-violet-600 to-emerald-400 p-0.5 shadow-md shadow-indigo-500/20 flex items-center justify-center">
              <div className="h-full w-full bg-white dark:bg-zinc-900 rounded-[10px] flex items-center justify-center">
                <span className="font-extrabold text-transparent bg-clip-text bg-gradient-to-tr from-indigo-600 to-violet-500 text-xl tracking-tight">
                  EB
                </span>
              </div>
            </div>
            <div>
              <div className="flex items-center gap-2">
                <h1 className="text-xl font-bold tracking-tight text-zinc-900 dark:text-zinc-50">
                  EventBeep
                </h1>
                <span className="inline-flex items-center gap-1 rounded-full bg-indigo-50 px-2 py-0.5 text-xs font-semibold text-indigo-700 ring-1 ring-inset ring-indigo-700/10 dark:bg-indigo-950/50 dark:text-indigo-300 dark:ring-indigo-400/20">
                  <SparklesIcon className="w-3 h-3 text-indigo-500" />
                  Radar
                </span>
              </div>
              <p className="text-xs text-zinc-500 dark:text-zinc-400">
                Automated Fresher Jobs & Internships Aggregator across India
              </p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 px-3 py-1.5 rounded-full border border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900 text-xs">
              <span className="relative flex h-2.5 w-2.5">
                <span
                  className={`animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${
                    isFallback ? "bg-amber-400" : "bg-emerald-400"
                  }`}
                />
                <span
                  className={`relative inline-flex rounded-full h-2.5 w-2.5 ${
                    isFallback ? "bg-amber-500" : "bg-emerald-500"
                  }`}
                />
              </span>
              <span className="font-medium text-zinc-700 dark:text-zinc-300">
                {isFallback
                  ? "Live Ingestion (Demo / Offline Mode)"
                  : "Go API & Gemini Ingestion Active"}
              </span>
              <span className="text-zinc-400 dark:text-zinc-600">|</span>
              <span className="text-zinc-500 dark:text-zinc-400 font-semibold">
                {totalJobs} Jobs Indexed
              </span>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
