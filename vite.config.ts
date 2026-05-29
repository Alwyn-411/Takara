import babel from '@rolldown/plugin-babel';
import react, { reactCompilerPreset } from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

// https://vite.dev/config/
export default defineConfig({
    plugins: [react(), babel({ presets: [reactCompilerPreset()] })],
    server: {
        proxy: {
            '/v1': 'http://localhost:8080',
            '/ping': 'http://localhost:8080',
        },
    },
});
