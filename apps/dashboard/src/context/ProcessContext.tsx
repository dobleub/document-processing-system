'use client'
import { createContext } from 'react'

export type Process = {
  id: string
  status: string
  error: string
  started_at: string
  estimated_completion: string
  files_processed: string[]
  files_to_process: string[]
  completed_at: string
}

export type ProcessContextType = {
  processes: Process[]
}

export const ProcessContext = createContext<Promise<ProcessContextType> | null>(null)

export const ProcessProvider = ({
  children,
  processPromise,
}: {
  children: React.ReactNode
  processPromise: Promise<ProcessContextType>
}) => {
  return <ProcessContext.Provider value={processPromise}>{children}</ProcessContext.Provider>
}
