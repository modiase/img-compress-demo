# Image Compression Demo

Web application for compressing images using DCT or SVD with real-time component visualization.

![Screenshot](/assets/screenshot.png)

## Features

- **DCT Compression**: Block-based (1-20 components)
- **SVD Compression**: Matrix decomposition (1-256 components)
- Interactive component slider with quality/size tradeoff

## Development

```bash
# Install dependencies
pushd frontend && pnpm install && popd
pushd backend && go mod download && popd

# Run both servers
pushd frontend && pnpm dev
```

Servers:

- Frontend: http://localhost:5173
- Backend: http://localhost:8080
