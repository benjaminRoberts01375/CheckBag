import { createRoot } from "react-dom/client";
import { StrictMode } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ContextProvider } from "./context";
import SignIn from "./screens/signin";
import SignUp from "./screens/signup";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<BrowserRouter>
			<ContextProvider>
				<Routes>
					<Route path="/" element={<h1>Hello, React!</h1>} />
					<Route path="/signin" element={<SignIn />} />
					<Route path="/signup" element={<SignUp />} />
				</Routes>
			</ContextProvider>
		</BrowserRouter>
	</StrictMode>,
);
