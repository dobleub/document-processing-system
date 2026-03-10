import React from "react";
import type { Metadata } from "next";
import ComponentCard from "@/components/common/ComponentCard";
import BasicTableMonitor from "@/components/monitor/BasicTableMonitor";
import { getAllProcess } from "@/lib/documentProcessing";

export const metadata: Metadata = {
  title:
    "Simple Monitor Dashboard | by TailAdmin",
  description: "This is Next.js Home used as a Monitor for Document Processing System",
};

export default async function Metrics() {
  const allProcess = await getAllProcess();

  return (
    <div className="w-full">
      <div className="space-y-6">
        <ComponentCard title="Simple Monitor Dashboard">
          <BasicTableMonitor allProcess={allProcess} />
        </ComponentCard>
      </div>
    </div>
  );
}
