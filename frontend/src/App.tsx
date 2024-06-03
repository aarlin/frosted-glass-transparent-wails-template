import DpsChart from "@/components/DpsChart"

function App() {
  return (
    <div className="isolate min-h-screen frosted-glass p-10 grid place-items-center mx-auto py-8">
      <div className="text-red-900 text-2xl font-bold flex flex-col items-center space-y-4">
        <DpsChart/>
      </div>
    </div>
  )
}

export default App
