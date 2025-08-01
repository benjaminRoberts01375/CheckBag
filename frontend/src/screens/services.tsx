import "../styles.css";
import servicesStyles from "./services.module.css";
import DashboardStyles from "./dashboard.module.css";
import { useList } from "../context-hook";
import ServiceEdit from "../components/service-edit";
import { useEffect } from "react";

const ServicesScreen = () => {
	const { services, requestServiceData } = useList();
	useEffect(() => {
		requestServiceData();
	}, []);

	return (
		<div id={servicesStyles["container"]}>
			<div className={DashboardStyles["graph-group"]}>
				<h2 className="header">Services</h2>
				<div id={servicesStyles["services"]}>
					{services.map(service => (
						<ServiceEdit service={service} key={service.clientID} />
					))}
					<ServiceEdit service={undefined} key={"new"} />
				</div>
			</div>
		</div>
	);
};

export default ServicesScreen;
