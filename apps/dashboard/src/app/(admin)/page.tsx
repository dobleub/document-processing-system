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
        <ComponentCard title="Simple Monitor Dashboard" suffixTitle={<small className="text-gray-500">(WebSocket)</small>}>
          <BasicTableMonitor allProcess={allProcess} />
          <div className="mb-4 text-sm text-gray-500">
            Note: *This is a simple monitor dashboard that displays the first 10 statuses of document processing in real-time using WebSocket.<br />
            You can see the list of processes, their IDs, and their current status. You can also start a new process using the <strong>Start New Process</strong> button.
          </div>
        </ComponentCard>
      </div>
    </div>
  );
}
