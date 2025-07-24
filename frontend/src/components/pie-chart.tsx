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
		<div>
			<h2 className={graphsStyles["header"]}>{title}</h2>
			<PieChart
				series={[
					{
						data: data,
					},
				]}
				width={200}
				height={200}
			/>
		</div>
	);
};

export default PieChartComponent;
