import { useEffect, useState } from 'react';
import { EventsOn } from '@wailsjs/runtime/runtime';
import { Button } from './button';
import { isEmpty } from '@/lib/utils';

const DpsChart = () => {
  // const [entities, setEntities] = useState({});
  const [logParsing, setLogParsing] = useState<boolean>(false);
  const [logs, setLogs] = useState<string[]>([]);
  const [logLines, setLogLines] = useState<number>(0);

  useEffect(() => {
    EventsOn("rcv:entities", (entities: any) => {
      if (!isEmpty(entities)) {
        setLogs(entities);
        console.log(entities)
      }
    })

    EventsOn("rcv:logLines", (logLines: number) => {
      console.log(logLines)
      setLogLines(logLines);
    })

    // Cleanup function to remove the event listener
    return () => {
      // window.backend.EventsOff('updateDPS');
    };
  }, []);

  const startDpsTracker = () => {
    window.go.main.App.StartDpsTracker();
    setLogParsing(true);
  }

  return (
    <div className="App">
      <header className="App-header flex flex-col justify-center items-center">
        <h1>DPS Calculator</h1>
        <div className="space-x-3">
          <Button onClick={startDpsTracker}>Start Tracking</Button>
        </div>
        <div id="dpsContainer">
          { logParsing ? 'Log is parsing' : 'Log is not being parsed'}
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
