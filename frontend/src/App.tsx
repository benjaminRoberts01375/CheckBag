import { createRoot } from "react-dom/client";
import { StrictMode } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import SignIn from "./screens/signin";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<BrowserRouter>
			<Routes>
				<Route path="/" element={<h1>Hello, React!</h1>} />
				<Route path="/signin" element={<SignIn />} />
			</Routes>
		</BrowserRouter>
	</StrictMode>,
);
