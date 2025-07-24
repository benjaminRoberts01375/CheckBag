import "../styles.css";
import GraphStyles from "./graphs.module.css";
import ResourceUsageData from "../types/resource-usage-data";

interface ResourceTableProps {
	data: ResourceUsageData[];
	title: string;
}

const ResourceTable = ({ data, title }: ResourceTableProps) => {
	return (
		<div id={GraphStyles["container"]}>
			<h2 className={GraphStyles["header"]}>{title}</h2>
			<table>
				<thead>
					<tr>
						<th>Service</th>
						<th>Resource</th>
						<th>Quantity</th>
					</tr>
				</thead>
				<tbody>
					{data.map(resourceUsage => {
						return (
							<tr>
								<td>{resourceUsage.service}</td>
								<td>{resourceUsage.resource}</td>
								<td>{resourceUsage.quantity}</td>
							</tr>
						);
					})}
				</tbody>
			</table>
		</div>
	);
};

export default ResourceTable;
