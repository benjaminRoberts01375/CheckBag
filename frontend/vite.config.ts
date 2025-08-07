import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
	plugins: [
		react({
			babel: {
				plugins: [["babel-plugin-react-compiler", {}]],
			},
		}),
	],
	build: {
		outDir: "dist",
	},
	server: {
		host: "0.0.0.0",
		port: 5173,
	},
});
