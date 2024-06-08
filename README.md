# Transparent / Frosted Glass Wails Template

## Technology

1. React
2. Tailwind
3. Wails

### Using the Template
```console
wails init -n project-name -t https://github.com/aarlin/frosted-glass-transparent-wails-template
```

```console
cd frontend
```

```console
npm install
```

### Configuration

Change app.go line 58 for frosted glass and transparent appearance:

Frosted Glass -> `WindowIsTranslucent: true`  
Transparent   -> `WindowIsTranslucent: false` 

## Live Development

To run in live development mode, run `wails dev` in the project directory. In another terminal, go into the `frontend`
directory and run `npm run dev`. The frontend dev server will run on http://localhost:34115. Connect to this in your
browser and connect to your application.

## Building

To build a redistributable, production mode package, use `wails build`.
