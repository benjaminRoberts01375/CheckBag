import "../styles.css";
import servicesStyles from "./services.module.css";
import { useList } from "../context-hook";
import { useEffect, useState } from "react";
import GraphData from "../types/graph-data";
import StackedBarChart from "../components/stacked-bar-chart";

const DashboardScreen = () => {
	const { services, requestServiceData, timescale } = useList();
	const [quantityData, setQuantityData] = useState<GraphData[]>([]);

	useEffect(() => {
		requestServiceData();
	}, []);

	useEffect(() => {
		var graphData: GraphData[] = [];

		for (const service of services.filter(service => service.enabled)) {
			const analyticsMap = service[timescale]; // Woo bracket notation
			analyticsMap.forEach((analytic, date) => {
				graphData.push(new GraphData(analytic.quantity, service.title, date));
			});
		}
		setQuantityData(graphData);
	}, [timescale, services]);

	return (
		<div id={servicesStyles["container"]}>
			<StackedBarChart graphData={quantityData} />
		</div>
	);
};

export default DashboardScreen;
