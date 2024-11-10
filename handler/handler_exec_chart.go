//go:build (all || most || chart) && !no_chart

package handler

// doExecChart executes a single query against the database, displaying its output as a chart.
func (h *Handler) doExecChart(ctx context.Context, w io.Writer, opt metacmd.Option, prefix, sqlstr string, qtyp bool, bind []interface{}) error {
	stdout, _, _ := h.l.Stdout(), h.l.Stderr(), h.l.Interactive()
	typ := env.TermGraphics()
	if !typ.Available() {
		return text.ErrGraphicsNotSupported
	}
	if _, ok := opt.Params["help"]; ok {
		fmt.Fprintln(stdout, text.ChartUsage)
		return nil
	}
	cfg, err := charts.ParseArgs(opt.Params)
	if err != nil {
		return err
	}
	start := time.Now()
	// query
	rows, err := h.DB().QueryContext(ctx, sqlstr, bind...)
	if err != nil {
		return err
	}
	// get cols
	cols, err := drivers.Columns(h.u, rows)
	if err != nil {
		return err
	}
	// process row(s)
	transposed := make([][]string, len(cols))
	clen, tfmt := len(cols), env.Vars().PrintTimeFormat()
	for rows.Next() {
		row, err := h.scan(rows, clen, tfmt)
		if err != nil {
			return err
		}
		for i := range row {
			transposed[i] = append(transposed[i], row[i])
		}
	}
	// display
	c, err := charts.MakeChart(cfg, cols, transposed)
	if err != nil {
		return err
	}
	data, err := c.ToEcharts()
	if err != nil {
		return err
	}
	echarts := echartsgoja.New(echartsgoja.WithWidthHeight(cfg.W, cfg.H))
	res, err := echarts.RenderOptions(ctx, data)
	if err != nil {
		return err
	}
	if cfg.File != "" {
		fmt.Println("writing to", cfg.File)
		return os.WriteFile(cfg.File, []byte(res), 0o644)
	}
	img, err := resvg.Render([]byte(res), resvg.WithBackground(cfg.Background))
	if err != nil {
		return err
	}
	if err := typ.Encode(stdout, img); err != nil {
		return err
	}
	if h.timing {
		d := time.Since(start)
		s := text.TimingDesc
		v := []interface{}{float64(d.Microseconds()) / 1000}
		if d > 1*time.Second {
			s += " (%v)"
			v = append(v, d.Round(1*time.Millisecond))
		}
		fmt.Fprintln(h.l.Stdout(), fmt.Sprintf(s, v...))
	}
	return nil
}
