import React from "react";
import { AlertCircleIcon, RefreshCwIcon } from "./Icons";

interface ErrorBannerProps {
  message?: string;
  onRetry: () => void;
}

export function ErrorBanner({
  message = "Failed to load jobs from the API. The backend may still be starting up.",
  onRetry,
}: ErrorBannerProps) {
  return (
    <div className="rounded-2xl border border-rose-200 bg-rose-50/70 dark:border-rose-900/50 dark:bg-rose-950/20 p-5 shadow-sm">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div className="flex items-start gap-3">
          <div className="mt-0.5 text-rose-600 dark:text-rose-400">
            <AlertCircleIcon className="h-5 w-5" />
          </div>
          <div>
            <h4 className="text-sm font-semibold text-rose-900 dark:text-rose-200">
              API Connection Warning
            </h4>
            <p className="text-xs text-rose-700 dark:text-rose-300/80 mt-0.5">
              {message}
            </p>
          </div>
        </div>
        <button
          onClick={onRetry}
          className="inline-flex items-center justify-center gap-1.5 self-start sm:self-auto rounded-xl bg-rose-600 px-3.5 py-2 text-xs font-semibold text-white shadow-sm hover:bg-rose-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 transition-colors"
        >
          <RefreshCwIcon className="h-3.5 w-3.5" />
          Retry Connection
        </button>
      </div>
    </div>
  );
}
