import "../styles.css";
import servicesStyles from "./services.module.css";
import { useList } from "../context-hook";
import ServiceEdit from "../components/service-edit";
import { useEffect } from "react";

// TODO: Temporary dashboard screen
const DashboardScreen = () => {
	const { services, requestServiceData } = useList();
	useEffect(() => {
		requestServiceData();
	}, []);

	return (
		<div id={servicesStyles["container"]}>
			<table id={servicesStyles["fancy-table"]}>
				<thead>
					<tr>
						<th>Name</th>
						<th>External Address</th>
						<th>Internal Address</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					{services.map(service => (
						<ServiceEdit service={service} key={service.clientID} />
					))}
					<ServiceEdit service={undefined} key={"new"} />
				</tbody>
			</table>
		</div>
	);
};

export default DashboardScreen;
