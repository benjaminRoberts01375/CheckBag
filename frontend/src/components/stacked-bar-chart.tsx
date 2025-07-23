import "../styles.css";
import GraphData from "../types/graph-data";

interface StackedBarChartProps {
	rawData: GraphData[];
}

const StackedBarChart = ({ rawData }: StackedBarChartProps) => {
	return (
		<div>
			<h1>Stacked Bar Chart</h1>
		</div>
	);
};

export default StackedBarChart;
