import { createRoot } from "react-dom/client";
import { StrictMode } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ContextProvider } from "./context";
import SignIn from "./screens/signin";
import SignUp from "./screens/signup";
import Dashboard from "./components/dashboard";
import DashboardScreen from "./screens/dashboard";
import ServicesScreen from "./screens/services";

createRoot(document.getElementById("root")!).render(
	<StrictMode>
		<BrowserRouter>
			<ContextProvider>
				<Routes>
					<Route path="/signin" element={<SignIn />} />
					<Route path="/signup" element={<SignUp />} />
					<Route path="/dashboard" element={<Dashboard />}>
						<Route index element={<Navigate to="/dashboard/home" replace />} />
						<Route path="home" element={<DashboardScreen />} />
						<Route path="services" element={<ServicesScreen />} />
						<Route path="*" element={<Navigate to="/dashboard/home" replace />} />
					</Route>
					<Route path="*" element={<Navigate to="/signin" replace />} />
				</Routes>
			</ContextProvider>
		</BrowserRouter>
	</StrictMode>,
);
