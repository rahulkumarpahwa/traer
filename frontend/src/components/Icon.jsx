export default function Icon({
  children,
  className = "h-5 w-5",
  fill = "none",
  stroke = "currentColor",
  strokeWidth = 1.8,
  ...props
}) {
  return (
    <svg
      viewBox="0 0 24 24"
      fill={fill}
      stroke={stroke}
      strokeWidth={strokeWidth}
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
      {...props}
    >
      {children}
    </svg>
  );
}
