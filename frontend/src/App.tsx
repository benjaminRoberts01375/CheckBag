import { createRoot } from "react-dom/client";
import { StrictMode } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<h1>Hello, world!</h1>} />
			</Routes>
		</BrowserRouter>
	</StrictMode>,
);
