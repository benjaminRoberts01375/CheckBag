import "../styles.css";
import DashboardStyles from "./dashboard.module.css";
import { Outlet } from "react-router-dom";
import Navbar from "./navbar";
import AnimatedBackground from "../components/animated-background";
import { useList } from "../context-hook";
import { useEffect, useState } from "react";

const Dashboard = () => {
	const { services } = useList();
	const [nodes, setNodes] = useState(services.length);

	useEffect(() => {
		console.log("Updating nodes to " + services.length);
		setNodes(services.length);
	}, [services]);

	return (
		<>
			<AnimatedBackground nodes={nodes} speed={0.5} />
			<div id={DashboardStyles["dashboard-container"]}>
				<Navbar />
				<div className={DashboardStyles["content-area"]}>
					<Outlet />
				</div>
			</div>
		</>
	);
};

export default Dashboard;
