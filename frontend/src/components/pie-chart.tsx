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
				<h2 className={graphsStyles["header"]}>{title}</h2>
				<PieChart
					series={[
						{
							data: data,
						},
					]}
					width={200}
					height={200}
					slotProps={{
						tooltip: {
							sx: {
								maxWidth: "700pt",
								minWidth: "50pt",
								whiteSpace: "nowrap",
								"& .MuiPaper-root": {
									maxWidth: "700px",
									minWidth: "100px",
								},
							},
						},
					}}
				/>
			</div>
		</ThemeProvider>
	);
};

export default PieChartComponent;
