const monthMap = new Map([
  ["Jan", "January"],
  ["Feb", "February"],
  ["Mar", "March"],
  ["Apr", "April"],
  ["May", "May"],
  ["Jun", "June"],
  ["Jul", "July"],
  ["Aug", "August"],
  ["Sep", "September"],
  ["Oct", "October"],
  ["Nov", "November"],
  ["Dec", "December"],
]);

// TODO rename
function myTooltipFormatter(p) {
  const color = p.color || "#666";
  const name = monthMap.get(p.name) || p.name;
  return (
    '<span style="color:white;font-size:14px;font-weight:medium;">' +
    name +
    ': <span style="color:' +
    color +
    ';font-weight:bold">' +
    "€" +
    p.value +
    "</span>" +
    "</span>"
  );
}
