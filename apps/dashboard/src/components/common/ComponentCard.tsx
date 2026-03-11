import React from "react";

interface ComponentCardProps {
  prefixTitle?: React.ReactNode; // Optional prefix title, can be a string or a React node
  title: string;
  suffixTitle?: React.ReactNode; // Optional suffix title, can be a string or a React node
  children: React.ReactNode;
  className?: string; // Additional custom classes for styling
  desc?: string; // Description text
}

const ComponentCard: React.FC<ComponentCardProps> = ({
  prefixTitle,
  title,
  suffixTitle,
  children,
  className = "",
  desc = "",
}) => {
  return (
    <div
      className={`rounded-2xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-white/[0.03] ${className}`}
    >
      {/* Card Header */}
      <div className="px-6 py-5">
        <h3 className="flex items-center text-base font-medium text-gray-800 dark:text-white/90">
          {prefixTitle && <span className="mr-2">{prefixTitle}</span>}
          {title}
          {suffixTitle && <span className="ml-2">{suffixTitle}</span>}
        </h3>
        {desc && (
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {desc}
          </p>
        )}
      </div>

      {/* Card Body */}
      <div className="p-4 border-t border-gray-100 dark:border-gray-800 sm:p-6">
        <div className="space-y-6">{children}</div>
      </div>
    </div>
  );
};

export default ComponentCard;
