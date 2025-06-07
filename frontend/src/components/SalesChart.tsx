interface SalesData {
  date: string;
  orders: number;
  purchases: number;
}

interface SalesChartProps {
  data: SalesData[];
  isLoading?: boolean;
}

export default function SalesChart({ data, isLoading = false }: SalesChartProps) {
  if (isLoading) {
    return (
      <div className="animate-pulse">
        <div className="h-64 bg-gray-200 rounded"></div>
      </div>
    );
  }

  if (!data || data.length === 0) {
    return (
      <div className="h-64 flex items-center justify-center text-gray-500">
        Нет данных для отображения
      </div>
    );
  }

  const maxValue = Math.max(...data.map(d => Math.max(d.orders, d.purchases)));
  const minValue = 0;
  const range = maxValue - minValue;
  
  return (
    <div className="h-full relative">
      <svg className="w-full h-full" viewBox="0 0 400 240" preserveAspectRatio="none">
        {/* Grid lines */}
        {[0, 1, 2, 3, 4].map(i => (
          <line
            key={i}
            x1="0"
            y1={i * 48}
            x2="400"
            y2={i * 48}
            stroke="#f3f4f6"
            strokeWidth="1"
          />
        ))}
        
        {/* Orders line (black) */}
        <polyline
          fill="none"
          stroke="#000000"
          strokeWidth="2"
          points={data
            .map((item, index) => {
              const x = (index / (data.length - 1)) * 400;
              const y = 240 - ((item.orders - minValue) / range) * 200;
              return `${x},${y}`;
            })
            .join(' ')}
        />
        
        {/* Purchases line (red) */}
        <polyline
          fill="none"
          stroke="#ef4444"
          strokeWidth="2"
          points={data
            .map((item, index) => {
              const x = (index / (data.length - 1)) * 400;
              const y = 240 - ((item.purchases - minValue) / range) * 200;
              return `${x},${y}`;
            })
            .join(' ')}
        />
        
        {/* Orders data points */}
        {data.map((item, index) => {
          const x = (index / (data.length - 1)) * 400;
          const y = 240 - ((item.orders - minValue) / range) * 200;
          return (
            <circle
              key={`orders-${index}`}
              cx={x}
              cy={y}
              r="3"
              fill="#000000"
            />
          );
        })}
        
        {/* Purchases data points */}
        {data.map((item, index) => {
          const x = (index / (data.length - 1)) * 400;
          const y = 240 - ((item.purchases - minValue) / range) * 200;
          return (
            <circle
              key={`purchases-${index}`}
              cx={x}
              cy={y}
              r="3"
              fill="#ef4444"
            />
          );
        })}
      </svg>
      
      {/* X-axis labels */}
      <div className="absolute bottom-0 left-0 right-0 flex justify-between text-xs text-gray-500 pt-2">
        {data.map((item, index) => (
          <span key={index} className="text-center">
            {item.date}
          </span>
        ))}
      </div>
    </div>
  );
} 