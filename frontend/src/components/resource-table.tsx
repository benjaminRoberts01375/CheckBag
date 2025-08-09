import "../styles.css";
import GraphStyles from "./graphs.module.css";
import ResourceTableStyles from "./resource-table.module.css";
import ResourceUsageData from "../types/resource-usage-data";

interface ResourceTableProps {
	data: ResourceUsageData[];
	title: string;
}

async function copy_path(path: string): Promise<void> {
	await navigator.clipboard.writeText(path);
}

const ResourceTable = ({ data, title }: ResourceTableProps) => {
	return (
		<div id={GraphStyles["container"]}>
			<h2 className="header">{title}</h2>
			{data.length === 0 ? (
				<p className={ResourceTableStyles["no-data"]}>No data to display</p>
			) : (
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
								<tr
									onClick={() =>
										copy_path(resourceUsage.service_url + "/" + resourceUsage.resource)
									}
									key={resourceUsage.service + resourceUsage.resource}
									className={ResourceTableStyles["clickable-row"]}
								>
									<td>
										<p className={ResourceTableStyles["service"]}>{resourceUsage.service}</p>
									</td>
									<td>
										<p className={ResourceTableStyles["resource"]}>{resourceUsage.resource}</p>
									</td>
									<td>
										<p>{resourceUsage.quantity}</p>
									</td>
								</tr>
							);
						})}
					</tbody>
				</table>
			)}
		</div>
	);
};

export default ResourceTable;
