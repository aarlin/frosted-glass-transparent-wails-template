// src/globals.d.ts
interface Window {
  startDpsTracker: () => void;
  go: {
    main: {
      App: {
        StartDpsTracker: () => void;
      };
    };
  };
}

// To ensure TypeScript recognizes the augmentation
declare var window: Window;
