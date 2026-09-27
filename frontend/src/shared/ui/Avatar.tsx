import { useMemo } from "react";

interface AvatarProps {
  username: string;
  avatarUrl?: string | null;
  size?: number;
  online?: boolean;
  className?: string;
  cacheBust?: number;
}

export function Avatar({
  username,
  avatarUrl,
  size = 40,
  online,
  className = "",
  cacheBust,
}: AvatarProps) {
  const initial = (username || "?")[0]?.toUpperCase() ?? "?";
  const bustedUrl = useMemo(
    () =>
      avatarUrl
        ? `${avatarUrl}${avatarUrl.includes("?") ? "&" : "?"}v=${cacheBust ?? 0}`
        : null,
    [avatarUrl, cacheBust],
  );

  return (
    <div
      className={`relative shrink-0 ${className}`}
      style={{ width: size, height: size }}
    >
      {bustedUrl ? (
        <img
          src={bustedUrl}
          alt={username}
          className="h-full w-full rounded-full object-cover"
          style={{ border: "1px solid var(--surface-border)" }}
        />
      ) : (
        <div
          className="flex h-full w-full items-center justify-center rounded-full font-display font-medium"
          style={{
            background:
              "linear-gradient(135deg,var(--theme-accent),var(--theme-accent-light))",
            color: "var(--theme-on-accent)",
            fontSize: size * 0.42,
          }}
        >
          {initial}
        </div>
      )}
      {online !== undefined && (
        <span
          className="absolute rounded-full"
          style={{
            width: Math.max(9, size * 0.26),
            height: Math.max(9, size * 0.26),
            right: -1,
            bottom: -1,
            background: online ? "#10b981" : "var(--text-tertiary)",
            border: "2px solid var(--surface-elevated)",
          }}
        />
      )}
    </div>
  );
}
