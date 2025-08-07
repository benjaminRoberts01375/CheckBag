import "../styles.css";
import DashboardStyles from "./dashboard.module.css";
import { Outlet } from "react-router-dom";
import Navbar from "./navbar";
import AnimatedBackground from "../components/animated-background";
import { useList } from "../context-hook";
import { useEffect, useState } from "react";
import { Fade as Hamburger } from "hamburger-react";

const Dashboard = () => {
	const { services } = useList();
	const [nodes, setNodes] = useState(services.length);
	const [isMobileView, setIsMobileView] = useState(window.innerWidth <= 884); // Initialize with current width
	const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

	useEffect(() => {
		console.log("Updating nodes to " + services.length);
		setNodes(services.length);
	}, [services.length]);

	// Check window size for mobile view
	useEffect(() => {
		const handleResize = () => {
			setIsMobileView(window.innerWidth <= 884);
		};

		window.addEventListener("resize", handleResize);

		// Cleanup
		return () => {
			window.removeEventListener("resize", handleResize);
		};
	}, []);

	return (
		<>
			<AnimatedBackground nodes={nodes} speed={0.5} />
			<div id={DashboardStyles["dashboard-container"]}>
				{isMobileView ? (
					<div id={DashboardStyles["hamburger-menu"]}>
						<Hamburger toggled={isMobileMenuOpen} toggle={setIsMobileMenuOpen} />
					</div>
				) : null}
				<Navbar
					isMobileView={isMobileView}
					isMobileMenuOpen={isMobileMenuOpen}
					setIsMobileMenuOpen={setIsMobileMenuOpen}
				/>
				<div className={DashboardStyles["content-area"]}>
					<Outlet />
				</div>
			</div>
		</>
	);
};

export default Dashboard;
