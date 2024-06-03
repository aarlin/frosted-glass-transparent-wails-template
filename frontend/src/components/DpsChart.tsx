import { useEffect, useState } from 'react';
import { EventsOn } from '@wailsjs/runtime/runtime';
import { Button } from './button';
import { isEmpty } from '@/lib/utils';

window.startDpsTracker = () => {
  // noinspection JSIgnoredPromiseFromCall
  window.go.main.App.StartDpsTracker();
}


const DpsChart = () => {
  // const [entities, setEntities] = useState({});
  const [logs, setLogs] = useState<string[]>([]);
  const [logLines, setLogLines] = useState<number>(0);

  useEffect(() => {
    EventsOn("rcv:entities", (entities: any) => {
      if (!isEmpty(entities))
      setLogs(entities);
    })

    EventsOn("rcv:logLines", (logLines: number) => {
      setLogLines(logLines);
    })

    // Cleanup function to remove the event listener
    return () => {
      // window.backend.EventsOff('updateDPS');
    };
  }, []);

  return (
    <div className="App">
      <header className="App-header flex flex-col justify-center items-center">
        <h1>DPS Calculator</h1>
        <div className="space-x-3">
          <Button onClick={window.startDpsTracker}>Start Tracking</Button>
          <Button onClick={window.startDpsTracker}>Stop Tracking</Button>
        </div>
        <div id="dpsContainer">
          {logLines}
          {
            logs.length > 0 && logs.map((log) => <div key={log}>{log}</div>)
          }
          {/* <h2>10s DPS:</h2>
          {Object.keys(entities).map((entityId) => (
            entities[entityId].dps10s !== 0 && (
              <div key={entityId}>
                {`ID: ${entityId}, Name: ${entities[entityId].name}, DPS10s: ${entities[entityId].dps10s.toFixed(2)}`}
              </div>
            )
          ))}
          <h2>60s DPS:</h2>
          {Object.keys(entities).map((entityId) => (
            entities[entityId].dps60s !== 0 && (
              <div key={entityId}>
                {`ID: ${entityId}, Name: ${entities[entityId].name}, DPS60s: ${entities[entityId].dps60s.toFixed(2)}`}
              </div>
            )
          ))}
          <h2>History:</h2>
          {Object.keys(entities).map((entityId) => (
            entities[entityId].dpsOnEnemy !== 0 && (
              <div key={entityId}>
                {`ID: ${entityId}, Name: ${entities[entityId].name}, DPS: ${entities[entityId].dpsOnEnemy.toFixed(2)}`}
              </div>
            )
          ))} */}
        </div>
      </header>
    </div>
  );
}

export default DpsChart;
