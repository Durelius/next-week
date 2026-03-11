import { useState, useEffect, useRef, useCallback } from "react";

// --- Types ---
interface Stop {
  stopId: string;
  stopName: string;
  stopLat: string;
  stopLon: string;
  locationType: string;
}

interface EdgeMetadata {
  tripId: string;
  routeId?: string;
  serviceId?: string;
  tripHeadsign?: string;
  tripShortName?: string;
  departure: number;
  arrival: number;
  sourceStopName: string;
  destStopName: string;
}

interface Edge {
  source: string;
  destination: string;
  metadata: EdgeMetadata;

}

// --- API ---
const API_BASE = "http://localhost:8080";

async function fetchStopsByName(name: string): Promise<Stop[]> {
  if (!name.trim()) return [];
  const res = await fetch(`${API_BASE}/stopbyname/${encodeURIComponent(name)}`);
  if (!res.ok) throw new Error("Failed to fetch stops");
  return res.json();
}

async function fetchRoute(fromId: string, toId: string, time: string): Promise<Edge[]> {
  const res = await fetch(`${API_BASE}/path/${fromId}/${toId}/${encodeURIComponent(time)}`);
  if (!res.ok) throw new Error("Failed to fetch route");
  return res.json();
}

function minutesToTime(minutes: number): string {
  const h = Math.floor(minutes / 60) % 24;
  const m = minutes % 60;
  return `${String(h).padStart(2, "0")}:${String(m).padStart(2, "0")}`;
}

function travelDuration(dep: number, arr: number): string {
  const diff = arr - dep;
  if (diff <= 0) return "–";
  const h = Math.floor(diff / 60);
  const m = diff % 60;
  return h > 0 ? `${h} h ${m} min` : `${m} min`;
}

// --- StopInput Component ---
function StopInput({
  label,
  icon,
  value,
  onSelect,
  placeholder,
}: {
  label: string;
  icon: React.ReactNode;
  value: Stop | null;
  onSelect: (stop: Stop) => void;
  placeholder: string;
}) {
  const [query, setQuery] = useState(value?.stopName ?? "");
  const [results, setResults] = useState<Stop[]>([]);
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [focused, setFocused] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (value) setQuery(value.stopName);
  }, [value]);

  const search = useCallback((q: string) => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (!q.trim()) { setResults([]); setOpen(false); return; }
    debounceRef.current = setTimeout(async () => {
      setLoading(true);
      try {
        const data = await fetchStopsByName(q);
        setResults(data ?? []);
        setOpen(true);
      } catch {
        setResults([]);
      } finally {
        setLoading(false);
      }
    }, 300);
  }, []);

  return (
    <div style={{ position: "relative", flex: 1 }}>
      <div style={{
        display: "flex", alignItems: "center", gap: 10,
        background: "#fff", borderRadius: 8,
        border: focused ? "2px solid #006CBF" : "2px solid #D0D7DE",
        padding: "10px 14px", transition: "border-color 0.15s",
        boxShadow: focused ? "0 0 0 3px rgba(0,108,191,0.12)" : "0 1px 3px rgba(0,0,0,0.07)"
      }}>
        <span style={{ color: "#006CBF", flexShrink: 0, display: "flex" }}>{icon}</span>
        <div style={{ flex: 1 }}>
          <div style={{ fontSize: 10, fontWeight: 700, letterSpacing: "0.08em", color: "#6B7A8D", textTransform: "uppercase", marginBottom: 2 }}>{label}</div>
          <input
            ref={inputRef}
            value={query}
            onChange={e => { setQuery(e.target.value); search(e.target.value); }}
            onFocus={() => { setFocused(true); if (query) setOpen(true); }}
            onBlur={() => { setFocused(false); setTimeout(() => setOpen(false), 150); }}
            placeholder={placeholder}
            style={{
              border: "none", outline: "none", width: "100%",
              fontSize: 16, fontFamily: "inherit", color: "#0F1923",
              background: "transparent", fontWeight: 500,
            }}
          />
        </div>
        {loading && (
          <div style={{ width: 16, height: 16, borderRadius: "50%", border: "2px solid #D0D7DE", borderTopColor: "#006CBF", animation: "spin 0.7s linear infinite", flexShrink: 0 }} />
        )}
      </div>

      {open && results.length > 0 && (
        <div style={{
          position: "absolute", top: "calc(100% + 4px)", left: 0, right: 0, zIndex: 100,
          background: "#fff", borderRadius: 8, boxShadow: "0 8px 32px rgba(0,0,0,0.14)",
          border: "1px solid #E3E8EF", overflow: "hidden", maxHeight: 260, overflowY: "auto"
        }}>
          {results.slice(0, 8).map((stop, i) => (
            <button
              key={stop.stopId + i}
              onMouseDown={() => { onSelect(stop); setQuery(stop.stopName); setOpen(false); }}
              style={{
                display: "flex", alignItems: "center", gap: 10, width: "100%",
                padding: "11px 14px", border: "none", background: "none",
                cursor: "pointer", textAlign: "left", fontFamily: "inherit",
                borderBottom: i < results.length - 1 ? "1px solid #F0F2F5" : "none",
                transition: "background 0.1s",
              }}
              onMouseEnter={e => (e.currentTarget.style.background = "#F0F6FF")}
              onMouseLeave={e => (e.currentTarget.style.background = "none")}
            >
              <span style={{ color: "#006CBF", fontSize: 13 }}>
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5"><circle cx="12" cy="12" r="3" /><path d="M12 2v3M12 19v3M2 12h3M19 12h3" /></svg>
              </span>
              <div>
                <div style={{ fontSize: 14, fontWeight: 600, color: "#0F1923" }}>{stop.stopName}</div>
                <div style={{ fontSize: 11, color: "#8A96A3", marginTop: 1 }}>Stop ID: {stop.stopId}</div>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

// --- Route Result ---
function RouteCard({ edges }: { edges: Edge[] }) {
  if (!edges || edges.length === 0) return (
    <div style={{ textAlign: "center", padding: "40px 0", color: "#8A96A3", fontStyle: "italic" }}>
      No route found between these stops.
    </div>
  );

  const totalDep = edges[0].metadata.departure;
  const totalArr = edges[edges.length - 1].metadata.arrival;

  // Group consecutive edges by tripId into "legs"
  const legs: Edge[][] = [];
  let current: Edge[] = [];
  for (const edge of edges) {
    if (current.length === 0 || current[0].metadata.tripId === edge.metadata.tripId) {
      current.push(edge);
    } else {
      legs.push(current);
      current = [edge];
    }
  }
  if (current.length > 0) legs.push(current);

  return (
    <div style={{ background: "#fff", borderRadius: 12, boxShadow: "0 2px 16px rgba(0,0,0,0.09)", overflow: "hidden", border: "1px solid #E3E8EF" }}>
      {/* Header */}
      <div style={{ background: "linear-gradient(135deg, #003F8A 0%, #006CBF 100%)", padding: "18px 24px", color: "#fff", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
        <div>
          <div style={{ fontSize: 13, opacity: 0.8, fontWeight: 500, marginBottom: 4 }}>Total journey</div>
          <div style={{ fontSize: 22, fontWeight: 700 }}>
            {minutesToTime(totalDep)} → {minutesToTime(totalArr)}
          </div>
        </div>
        <div style={{ textAlign: "right" }}>
          <div style={{ fontSize: 13, opacity: 0.8 }}>Duration</div>
          <div style={{ fontSize: 20, fontWeight: 700 }}>{travelDuration(totalDep, totalArr)}</div>
        </div>
      </div>

      {/* Legs */}
      <div style={{ padding: "0 24px 20px" }}>
        {legs.map((leg, li) => {
          const legFrom = leg[0].source;
          const legTo = leg[leg.length - 1].destination;
          const legDep = leg[0].metadata.departure;
          const legArr = leg[leg.length - 1].metadata.arrival;
          const meta = leg[0].metadata;

          return (
            <div key={li}>
              {/* Transfer indicator */}
              {li > 0 && (
                <div style={{ display: "flex", alignItems: "center", gap: 8, padding: "8px 0", color: "#8A96A3", fontSize: 12 }}>
                  <div style={{ width: 2, height: 28, background: "#E3E8EF", marginLeft: 11 }} />
                  <span style={{ marginLeft: 4 }}>Transfer</span>
                </div>
              )}

              {/* Leg row */}
              <div style={{ display: "flex", gap: 16, paddingTop: li === 0 ? 20 : 0 }}>
                {/* Timeline */}
                <div style={{ display: "flex", flexDirection: "column", alignItems: "center", flexShrink: 0, width: 24 }}>
                  <div style={{ width: 12, height: 12, borderRadius: "50%", background: "#006CBF", border: "2px solid #fff", boxShadow: "0 0 0 2px #006CBF", zIndex: 1 }} />
                  <div style={{ flex: 1, width: 2, background: "#D6E4F3", minHeight: 40 }} />
                  <div style={{ width: 12, height: 12, borderRadius: "50%", background: "#003F8A", border: "2px solid #fff", boxShadow: "0 0 0 2px #003F8A", zIndex: 1 }} />
                </div>

                {/* Content */}
                <div style={{ flex: 1 }}>
                  {/* Departure */}
                  <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
                    <div>
                      <div style={{ fontSize: 15, fontWeight: 700, color: "#0F1923" }}>{legFrom}</div>
                      <div style={{ fontSize: 12, color: "#8A96A3", marginTop: 2 }}>
                        {meta.tripHeadsign && <span>Towards <strong>{meta.tripHeadsign}</strong> · </span>}
                        {meta.tripShortName && <span>Line {meta.tripShortName}</span>}
                      </div>
                    </div>
                    <div style={{ fontSize: 16, fontWeight: 700, color: "#006CBF", textAlign: "right", flexShrink: 0, marginLeft: 12 }}>
                      {minutesToTime(legDep)}
                    </div>
                  </div>

                  {/* Trip badge */}
                  <div style={{ margin: "10px 0", display: "inline-flex", alignItems: "center", gap: 6, background: "#F0F6FF", borderRadius: 6, padding: "5px 10px", border: "1px solid #C8DEFF" }}>
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#006CBF" strokeWidth="2.5">
                      <rect x="2" y="7" width="20" height="13" rx="2" /><path d="M16 7V5a2 2 0 00-2-2h-4a2 2 0 00-2 2v2" /><line x1="12" y1="12" x2="12" y2="16" /><line x1="10" y1="14" x2="14" y2="14" />
                    </svg>
                    <span style={{ fontSize: 12, color: "#006CBF", fontWeight: 600 }}>{leg.length} stop{leg.length !== 1 ? "s" : ""}</span>
                    <span style={{ fontSize: 11, color: "#5A8ECC" }}>· {travelDuration(legDep, legArr)}</span>
                  </div>

                  {/* Stops expand */}
                  <details style={{ marginBottom: 8 }}>
                    <summary style={{ cursor: "pointer", fontSize: 12, color: "#5A8ECC", userSelect: "none", listStyle: "none", display: "flex", alignItems: "center", gap: 4 }}>
                      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5"><path d="M6 9l6 6 6-6" /></svg>
                      View intermediate stops
                    </summary>
                    <div style={{ marginTop: 6, paddingLeft: 8, borderLeft: "2px solid #D6E4F3" }}>
                      {leg.map((e, ei) => (
                        <div key={ei} style={{ display: "flex", justifyContent: "space-between", padding: "3px 0", fontSize: 12, color: "#4A5568" }}>
                          <span>{e.metadata.destStopName}</span>
                          <span style={{ color: "#8A96A3" }}>{minutesToTime(e.metadata.departure)}</span>
                        </div>
                      ))}
                    </div>
                  </details>

                  {/* Arrival */}
                  <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                    <div style={{ fontSize: 15, fontWeight: 700, color: "#0F1923" }}>{legTo}</div>
                    <div style={{ fontSize: 16, fontWeight: 700, color: "#003F8A" }}>{minutesToTime(legArr)}</div>
                  </div>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

// --- Main App ---
export default function App() {
  const [from, setFrom] = useState<Stop | null>(null);
  const [to, setTo] = useState<Stop | null>(null);
  const [time, setTime] = useState(() => {
    const now = new Date();
    return now.toTimeString().slice(0, 5);
  });
  const [route, setRoute] = useState<Edge[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const canSearch = from && to;

  const handleSearch = async () => {
    if (!from || !to) return;
    setLoading(true);
    setError(null);
    setRoute(null);
    try {
      const data = await fetchRoute(from.stopId, to.stopId, time);
      setRoute(data);
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Could not fetch route");
    } finally {
      setLoading(false);
    }
  };

  const swap = () => {
    const tmp = from;
    setFrom(to);
    setTo(tmp);
  };

  return (
    <div style={{ minHeight: "100vh", background: "#F0F4F8", fontFamily: "'Figtree', 'Helvetica Neue', Helvetica, Arial, sans-serif" }}>
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Figtree:wght@400;500;600;700;800&display=swap');
        * { box-sizing: border-box; margin: 0; padding: 0; }
        @keyframes spin { to { transform: rotate(360deg); } }
        @keyframes fadeUp { from { opacity: 0; transform: translateY(16px); } to { opacity: 1; transform: translateY(0); } }
        @keyframes shimmer { 0% { background-position: -400px 0; } 100% { background-position: 400px 0; } }
        details > summary::-webkit-details-marker { display: none; }
        ::-webkit-scrollbar { width: 5px; } ::-webkit-scrollbar-thumb { background: #C0CCDA; border-radius: 4px; }
      `}</style>

      {/* Header */}
      <header style={{ background: "#003F8A", color: "#fff", boxShadow: "0 2px 12px rgba(0,0,0,0.25)" }}>
        <div style={{ maxWidth: 900, margin: "0 auto", padding: "0 20px", display: "flex", alignItems: "center", justifyContent: "space-between", height: 60 }}>
          {/* SL Logo */}
          <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
            <div style={{
              width: 44, height: 44, borderRadius: 10,
              background: "#fff", display: "flex", alignItems: "center", justifyContent: "center",
            }}>
              <svg width="32" height="24" viewBox="0 0 32 24" fill="none">
                <text x="0" y="20" fontFamily="'Figtree', Arial, sans-serif" fontWeight="800" fontSize="22" fill="#003F8A" letterSpacing="-1">SL</text>
              </svg>
            </div>
            <div>
              <div style={{ fontWeight: 800, fontSize: 17, letterSpacing: "-0.02em" }}>Stockholms Lokaltrafik</div>
              <div style={{ fontSize: 11, opacity: 0.65, fontWeight: 500 }}>Journey Planner</div>
            </div>
          </div>

          <nav style={{ display: "flex", gap: 4 }}>
            {["Timetables", "Tickets", "Maps", "About"].map(item => (
              <button key={item} style={{ background: "none", border: "none", color: "rgba(255,255,255,0.75)", padding: "6px 12px", borderRadius: 6, cursor: "pointer", fontSize: 13, fontWeight: 500, fontFamily: "inherit", transition: "all 0.15s" }}
                onMouseEnter={e => { (e.currentTarget as HTMLButtonElement).style.background = "rgba(255,255,255,0.12)"; (e.currentTarget as HTMLButtonElement).style.color = "#fff"; }}
                onMouseLeave={e => { (e.currentTarget as HTMLButtonElement).style.background = "none"; (e.currentTarget as HTMLButtonElement).style.color = "rgba(255,255,255,0.75)"; }}
              >{item}</button>
            ))}
          </nav>
        </div>
      </header>

      {/* Hero */}
      <div style={{ background: "linear-gradient(160deg, #003F8A 0%, #0062B8 55%, #1A80D4 100%)", padding: "48px 20px 80px" }}>
        <div style={{ maxWidth: 720, margin: "0 auto", textAlign: "center" }}>
          <h1 style={{ color: "#fff", fontSize: 36, fontWeight: 800, letterSpacing: "-0.03em", marginBottom: 8, lineHeight: 1.1 }}>
            Where are you going?
          </h1>
          <p style={{ color: "rgba(255,255,255,0.7)", fontSize: 16, marginBottom: 36 }}>
            Find the fastest route across Stockholm's public transit network
          </p>

          {/* Search Card */}
          <div style={{ background: "#fff", borderRadius: 16, padding: 20, boxShadow: "0 16px 48px rgba(0,0,0,0.22)", position: "relative" }}>
            <div style={{ display: "flex", gap: 12, alignItems: "stretch", flexWrap: "wrap" }}>
              <StopInput
                label="From"
                icon={<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5"><circle cx="12" cy="12" r="4" /><circle cx="12" cy="12" r="9" strokeDasharray="3 2" /></svg>}
                value={from}
                onSelect={setFrom}
                placeholder="Departure stop…"
              />

              {/* Swap */}
              <button onClick={swap} style={{
                background: "#F0F6FF", border: "2px solid #C8DEFF", borderRadius: 8, width: 44, flexShrink: 0,
                cursor: "pointer", display: "flex", alignItems: "center", justifyContent: "center", color: "#006CBF",
                transition: "all 0.15s", alignSelf: "center"
              }}
                onMouseEnter={e => { (e.currentTarget as HTMLButtonElement).style.background = "#006CBF"; (e.currentTarget as HTMLButtonElement).style.color = "#fff"; }}
                onMouseLeave={e => { (e.currentTarget as HTMLButtonElement).style.background = "#F0F6FF"; (e.currentTarget as HTMLButtonElement).style.color = "#006CBF"; }}
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                  <path d="M7 16V4m0 0L3 8m4-4l4 4M17 8v12m0 0l4-4m-4 4l-4-4" />
                </svg>
              </button>

              <StopInput
                label="To"
                icon={<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5"><path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z" /><circle cx="12" cy="9" r="2.5" /></svg>}
                value={to}
                onSelect={setTo}
                placeholder="Destination stop…"
              />
            </div>

            {/* Second row */}
            <div style={{ display: "flex", gap: 12, marginTop: 12, alignItems: "center" }}>
              <div style={{ display: "flex", alignItems: "center", gap: 8, background: "#F8FAFC", border: "2px solid #D0D7DE", borderRadius: 8, padding: "10px 14px", flex: "0 0 auto" }}>
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#6B7A8D" strokeWidth="2.5">
                  <circle cx="12" cy="12" r="9" /><path d="M12 6v6l4 2" />
                </svg>
                <input
                  type="time"
                  value={time}
                  onChange={e => setTime(e.target.value)}
                  style={{ border: "none", outline: "none", background: "transparent", fontSize: 14, fontFamily: "inherit", fontWeight: 600, color: "#0F1923", width: 80 }}
                />
              </div>

              <button
                onClick={handleSearch}
                disabled={!canSearch || loading}
                style={{
                  flex: 1, padding: "12px 24px", borderRadius: 8, border: "none",
                  background: canSearch ? "linear-gradient(135deg, #0062B8 0%, #1A80D4 100%)" : "#C0CCDA",
                  color: "#fff", fontSize: 16, fontWeight: 700, cursor: canSearch ? "pointer" : "not-allowed",
                  fontFamily: "inherit", display: "flex", alignItems: "center", justifyContent: "center", gap: 8,
                  boxShadow: canSearch ? "0 4px 14px rgba(0,108,191,0.4)" : "none",
                  transition: "all 0.2s", letterSpacing: "-0.01em",
                }}
                onMouseEnter={e => { if (canSearch) (e.currentTarget as HTMLButtonElement).style.transform = "translateY(-1px)"; }}
                onMouseLeave={e => { (e.currentTarget as HTMLButtonElement).style.transform = "none"; }}
              >
                {loading ? (
                  <div style={{ width: 18, height: 18, borderRadius: "50%", border: "2.5px solid rgba(255,255,255,0.3)", borderTopColor: "#fff", animation: "spin 0.7s linear infinite" }} />
                ) : (
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                    <circle cx="11" cy="11" r="7" /><path d="M21 21l-4.35-4.35" />
                  </svg>
                )}
                {loading ? "Finding route…" : "Search journey"}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Results */}
      <div style={{ maxWidth: 720, margin: "-32px auto 60px", padding: "0 20px" }}>
        {error && (
          <div style={{ background: "#FFF0F0", border: "1px solid #FFC5C5", borderRadius: 10, padding: "14px 18px", color: "#B91C1C", fontSize: 14, fontWeight: 500, marginBottom: 20, display: "flex", gap: 10, alignItems: "center", animation: "fadeUp 0.3s ease" }}>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5"><circle cx="12" cy="12" r="9" /><line x1="12" y1="8" x2="12" y2="12" /><line x1="12" y1="16" x2="12.01" y2="16" /></svg>
            {error}
          </div>
        )}

        {route !== null && (
          <div style={{ animation: "fadeUp 0.35s ease" }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 14 }}>
              <h2 style={{ fontSize: 18, fontWeight: 700, color: "#0F1923", letterSpacing: "-0.02em" }}>
                Journey result
              </h2>
              <span style={{ fontSize: 13, color: "#8A96A3" }}>{from?.stopName} → {to?.stopName}</span>
            </div>
            <RouteCard edges={route} />
          </div>
        )}

        {!route && !loading && !error && (
          <div style={{ textAlign: "center", paddingTop: 24 }}>
            <div style={{ display: "flex", gap: 16, justifyContent: "center", flexWrap: "wrap" }}>
              {[
                { icon: "🚇", label: "Tunnelbana" },
                { icon: "🚌", label: "Buses" },
                { icon: "🚋", label: "Trams" },
                { icon: "🚢", label: "Ferries" },
              ].map(item => (
                <div key={item.label} style={{ background: "#fff", borderRadius: 10, padding: "14px 20px", display: "flex", flexDirection: "column", alignItems: "center", gap: 6, boxShadow: "0 1px 6px rgba(0,0,0,0.07)", border: "1px solid #E3E8EF", minWidth: 90 }}>
                  <span style={{ fontSize: 26 }}>{item.icon}</span>
                  <span style={{ fontSize: 12, fontWeight: 600, color: "#4A5568" }}>{item.label}</span>
                </div>
              ))}
            </div>
            <p style={{ marginTop: 20, color: "#8A96A3", fontSize: 14 }}>Enter a departure and destination above to plan your journey</p>
          </div>
        )}
      </div>

      {/* Footer */}
      <footer style={{ background: "#1A2B45", color: "rgba(255,255,255,0.5)", padding: "24px 20px", textAlign: "center", fontSize: 12 }}>
        <p>© {new Date().getFullYear()} Stockholms Lokaltrafik AB · This is a local dev build</p>
      </footer>
    </div>
  );
}
