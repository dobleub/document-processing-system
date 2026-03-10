"use client";

import React from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableFooter,
  TableHeader,
  TableRow,
} from "../ui/table";

import Badge from "../ui/badge/Badge";

import { ProcessContextType } from "@/context/ProcessContext";
import { useGetAllProcessWebSocket } from "@/lib/documentProcessingWS";
import Button from "../ui/button/Button";
import { startProcess, stopProcess } from "@/lib/documentProcessing";

export default function BasicTableMonitor({ allProcess }: { allProcess: ProcessContextType }) {

  const { processes, isConnected, readyState, connectionStatus } = useGetAllProcessWebSocket();

  return (
    <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-white/[0.05] dark:bg-white/[0.03]">
      <div className="max-w-full overflow-x-auto">
        <div className="min-w-[1102px]">
          <div className="flex items-center justify-between px-5 py-4">
            <div className="flex items-center gap-2">
              {/* Button start process */}
              <Button size="sm" onClick={() => startProcess("front")}>Start Process</Button>
              <Badge size="md" color={isConnected ? "success" : "error"}>
                {connectionStatus}
              </Badge>
            </div>
          </div>
        </div>
        <div className="min-w-[1102px]">
          <Table>
            {/* Table Header */}
            <TableHeader className="border-b border-gray-100 dark:border-white/[0.05]">
              <TableRow>
                <TableCell
                  isHeader
                  className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400"
                >
                  ID
                </TableCell>
                <TableCell
                  isHeader
                  className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400"
                >
                  Status
                </TableCell>
                <TableCell
                  isHeader
                  className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400"
                >
                  Started At
                </TableCell>
                <TableCell
                  isHeader
                  className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400"
                >
                  Estimated Completion
                </TableCell>
                <TableCell
                  isHeader
                  className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400"
                >
                  Completed At
                </TableCell>
                <TableCell
                  isHeader
                  className="px-5 py-3 font-medium text-gray-500 text-start text-theme-xs dark:text-gray-400"
                >
                  Actions
                </TableCell>
              </TableRow>
            </TableHeader>

            {/* Table Body */}
            <TableBody className="divide-y divide-gray-100 dark:divide-white/[0.05]">
              {processes.map((p) => (
                <TableRow key={p.id}>
                  <TableCell className="px-5 py-4 sm:px-6 text-start">
                    <div className="flex items-center gap-3">
                      <div>
                        <span className="block font-medium text-gray-800 text-theme-sm dark:text-white/90">
                          {p.id}
                        </span>
                      </div>
                    </div>
                  </TableCell>
                  <TableCell className="px-4 py-3 text-gray-500 text-start text-theme-sm dark:text-gray-400">
                    {/*
                      ● COMPLETED - success
                      ● PENDING - warning
                      ● RUNNING - primary
                      ● PAUSED - info
                      ● FAILED - error
                      ● STOPPED - error

                      type BadgeColor =
                      | "primary"
                      | "success"
                      | "error"
                      | "warning"
                      | "info"
                      | "light"
                      | "dark";
                    */}
                    <Badge
                      size="sm"
                      color={
                        p.status === "COMPLETED"
                          ? "success"
                          : p.status === "PENDING"
                          ? "warning"
                          : p.status === "RUNNING"
                          ? "primary"
                          : p.status === "PAUSED"
                          ? "info"
                          : p.status === "FAILED"
                          ? "error"
                          : "error"
                      }
                    >
                      {p.status}
                    </Badge>
                  </TableCell>
                  <TableCell className="px-4 py-3 text-gray-500 text-start text-theme-sm dark:text-gray-400">
                    {p.started_at ? new Date(p.started_at).toLocaleString() : "N/A"}
                  </TableCell>
                  <TableCell className="px-4 py-3 text-gray-500 text-start text-theme-sm dark:text-gray-400">
                    <div className="flex -space-x-2">
                      {p.estimated_completion}
                    </div>
                  </TableCell>
                  <TableCell className="px-4 py-3 text-gray-500 text-theme-sm dark:text-gray-400">
                    {p.completed_at ? new Date(p.completed_at).toLocaleString() : "N/A"}
                  </TableCell>
                  <TableCell className="px-4 py-3 text-gray-500 text-start text-theme-sm dark:text-gray-400">
                    <Button
                      className="bg-orange-500 hover:bg-orange-600 focus:outline-2 focus:outline-offset-2 focus:outline-orange-500 active:bg-orange-700"
                      size="sm"
                      onClick={() => stopProcess(p.id, "front")}
                      disabled={p.status !== "RUNNING" && p.status !== "PENDING"}
                    >
                      Stop
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </div>
    </div>
  );
}
