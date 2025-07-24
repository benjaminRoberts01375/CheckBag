import "../styles.css";
import DashboardStyles from "./dashboard.module.css";
import { Outlet } from "react-router-dom";
import Navbar from "./navbar";

const Dashboard = () => {
	return (
		<div id={DashboardStyles["dashboard-container"]}>
			<Navbar />
			<div className={DashboardStyles["content-area"]}>
				<Outlet />
			</div>
		</div>
	);
};

export default Dashboard;
