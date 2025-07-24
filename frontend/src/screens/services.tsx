import "../styles.css";
import servicesStyles from "./services.module.css";
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
			<table id={servicesStyles["fancy-table"]}>
				<thead>
					<tr>
						<th>
							<p>Name</p>
						</th>
						<th>
							<p>External Address</p>
						</th>
						<th>
							<p>Internal Address</p>
						</th>
						<th>
							<p>Actions</p>
						</th>
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

export default ServicesScreen;
