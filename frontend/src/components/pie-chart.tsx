import "../styles.css";
import graphsStyles from "./graphs.module.css";
import ChartData from "../types/chart-data";
import { PieChart } from "@mui/x-charts/PieChart";
import { ThemeProvider, Theme } from "@mui/material/styles";

interface PieChartProps {
	data: ChartData[];
	title: string;
	theme: Theme;
}

const PieChartComponent = ({ data, title, theme }: PieChartProps) => {
	return (
		<ThemeProvider theme={theme}>
			<div id={graphsStyles["container"]}>
				<h2 className="header">{title}</h2>
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
		</ThemeProvider>
	);
};

export default PieChartComponent;
