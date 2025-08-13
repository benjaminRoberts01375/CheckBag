import "../styles.css";
import GraphStyles from "./graphs.module.css";
import ResourceTableStyles from "./resource-table.module.css";
import ResourceUsageData from "../types/resource-usage-data";
import { IoIosDownload } from "react-icons/io";

interface ResourceTableProps {
	data: ResourceUsageData[];
	title: string;
}

async function copy_path(path: string): Promise<void> {
	if (navigator.clipboard && window.isSecureContext) {
		try {
			await navigator.clipboard.writeText(path);
			return;
		} catch (err) {
			console.log("Clipboard API failed:", err);
		}
	}
	try {
		const textArea = document.createElement("textarea");
		textArea.value = path;
		textArea.style.position = "fixed";
		textArea.style.left = "-999999px";
		textArea.style.top = "-999999px";
		document.body.appendChild(textArea);
		textArea.focus();
		textArea.select();

		document.execCommand("copy"); // I wish we had another way to do this >:(
		document.body.removeChild(textArea);
		return;
	} catch (err) {
		console.error("Fallback copy failed:", err);
		return;
	}
}

const ResourceTable = ({ data, title }: ResourceTableProps) => {
	const exportToCSV = () => {
		if (data.length === 0) {
			console.warn("No data to export");
			return;
		}

		// Create CSV headers
		const headers = ["Service", "Resource", "Quantity"];

		// Convert data to CSV rows
		const csvRows = data.map(item => [
			`"${item.service}"`,
			`"${item.resource}"`,
			item.quantity.toString(),
		]);

		// Combine headers and rows
		const csvContent = [headers.join(","), ...csvRows.map(row => row.join(","))].join("\n");

		// Create and download the file
		const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
		const link = document.createElement("a");
		const url = URL.createObjectURL(blob);

		link.setAttribute("href", url);
		link.setAttribute("download", `${title.toLowerCase().replace(/\s+/g, "-")}-data.csv`);
		link.style.visibility = "hidden";
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
		URL.revokeObjectURL(url);
	};

	return (
		<div id={GraphStyles["container"]}>
			<div id={GraphStyles["header"]}>
				<h2>{title}</h2>
				<button
					onClick={exportToCSV}
					id={GraphStyles["save-button"]}
					className="secondary"
					title="Download CSV"
					disabled={data.length === 0}
				>
					<IoIosDownload />
				</button>
			</div>
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
