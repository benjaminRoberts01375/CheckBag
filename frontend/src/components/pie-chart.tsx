import "../styles.css";
import graphsStyles from "./graphs.module.css";
import ChartData from "../types/chart-data";
import { PieChart } from "@mui/x-charts/PieChart";

interface PieChartProps {
	data: ChartData[];
	title: string;
}

const PieChartComponent = ({ data, title }: PieChartProps) => {
	return (
		<div id={graphsStyles["container"]}>
			<h2 className="header">{title}</h2>
			<PieChart
				series={[
					{
						data: data,
					},
				]}
				width={170}
				height={170}
				margin={{ top: 0, right: 0, left: 0, bottom: 0 }}
			/>
		</div>
	);
};

export default PieChartComponent;
