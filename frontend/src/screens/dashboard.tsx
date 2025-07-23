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
			var analyticsMap = service[timescale]; // Woo bracket notation
			var timeStepQuantity = 0;
			var rollback: (step: number) => Date;

			switch (timescale) {
				case "hour":
					rollback = (step: number) => {
						const now = new Date();
						const currentMinute = now.getMinutes();
						now.setMinutes(currentMinute + step);
						now.setSeconds(0);
						now.setMilliseconds(0);
						return now;
					};
					timeStepQuantity = 60;
					break;
				case "day":
					rollback = (step: number) => {
						const now = new Date();
						const currentHour = now.getHours();
						now.setHours(currentHour + step);
						now.setMinutes(0);
						now.setSeconds(0);
						now.setMilliseconds(0);
						return now;
					};
					timeStepQuantity = 24;
					break;
				case "month":
					rollback = (step: number): Date => {
						var now = new Date();
						const currentUTCDay = now.getUTCDate();
						now.setUTCDate(currentUTCDay + step);
						now.setUTCHours(0);
						now.setUTCMinutes(0);
						now.setUTCSeconds(0);
						now.setUTCMilliseconds(0);
						return now;
					};
					timeStepQuantity = 30;
					break;
				case "year":
					rollback = (step: number) => {
						const now = new Date();
						const currentMonth = now.getMonth();
						now.setUTCMonth(currentMonth + step);
						now.setUTCDate(1);
						now.setUTCHours(0);
						now.setUTCMinutes(0);
						now.setUTCSeconds(0);
						now.setUTCMilliseconds(0);
						return now;
					};
					timeStepQuantity = 12;
					break;
			}

			for (let i = 0; i < timeStepQuantity; i++) {
				const date = rollback(-i);
				var analytic = analyticsMap.get(date.toISOString());
				graphData.push(new GraphData(analytic?.quantity ?? 0, service.title, date));
			}
		}
		setQuantityData(graphData);
	}, [timescale, services]);

	return (
		<div id={servicesStyles["container"]}>
			<StackedBarChart
				graphData={quantityData}
				timescale={timescale}
				yAxisLabel="Query Quantity"
				title="Query Quantity Per Service"
			/>
		</div>
	);
};

export default DashboardScreen;
