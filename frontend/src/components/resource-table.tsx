import "../styles.css";
import GraphStyles from "./graphs.module.css";
import ResourceTableStyles from "./resource-table.module.css";
import ResourceUsageData from "../types/resource-usage-data";

interface ResourceTableProps {
	data: ResourceUsageData[];
	title: string;
}

const ResourceTable = ({ data, title }: ResourceTableProps) => {
	return (
		<div id={GraphStyles["container"]}>
			<h2 className={GraphStyles["header"]}>{title}</h2>
			<table className={ResourceTableStyles["styled-table"]}>
				<thead>
					<tr>
						<th>
							<p>Service</p>
						</th>
						<th>
							<p>Resource</p>
						</th>
						<th>
							<p>Quantity</p>
						</th>
					</tr>
				</thead>
				<tbody>
					{data.map(resourceUsage => {
						return (
							<tr>
								<td>
									<p>{resourceUsage.service}</p>
								</td>
								<td>
									<p>{resourceUsage.resource}</p>
								</td>
								<td>
									<p>{resourceUsage.quantity}</p>
								</td>
							</tr>
						);
					})}
				</tbody>
			</table>
		</div>
	);
};

export default ResourceTable;
