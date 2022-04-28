package sysgapp

import (
	"sync"
	"unicode"

	V "github.com/gabe-lee/genvecs"
)

type VertexMode uint8

const (
	Pixels        VertexMode = iota // Each vertex is an independant pixel
	Lines                           // Each pair of vertices forms an independant line
	LineStrip                       // Each vertex forms a continuous line with the vertex following it
	LineLoop                        // Each vertex forms a continuous line with the vertex following it, last vertex connects back to first
	Triangles                       // Every 3 vertices form independant triangles
	TriangleStrip                   // Every vertex forms a triangle using the 2 following it with alternating windings
	TriangleFan                     // Every vertex uses the one following it and the very first vertex to form a triangle
) // Draw Modes

type ShaderType uint8

// Shader Type Codes
const (
	VertexShader ShaderType = iota
	FragmentShader
	GeometryShader
	ComputeShader
)

type Shader struct {
	sType ShaderType
	code  string
}

func NewShader(sType ShaderType, code string) *Shader {
	return &Shader{
		sType: sType,
		code:  code,
	}
}

type ImageType byte

// Image Type Codes
const (
	PNG ImageType = iota
	BMP
	WEBP
)

type Texture struct {
	data    []byte
	imgType ImageType
	size    V.F32Vec2
	mipMaps int32
	ID      uint32
	Unit    uint32
}

func NewTexture(data []byte, imgType ImageType, size V.F32Vec2, mipMaps int32) *Texture {
	return &Texture{
		data:    data,
		imgType: imgType,
		size:    size,
		mipMaps: mipMaps,
	}
}

// var PlanetSweeperTex = NewTexture(PlanetSweeperTexWEBP, WEBP, V.F32Vec2{512, 1024}, 0)

type TextureIndex int
type SurfaceIndex int

const (
	MainSurface           SurfaceIndex = iota
	MainTexture           TextureIndex = iota
	MapSurface            SurfaceIndex = iota
	SpriteAssemblySurface SurfaceIndex = iota
	//
	MapTexture            = TextureIndex(MapSurface)
	SpriteAssemblyTexture = TextureIndex(SpriteAssemblySurface)
) // Texture and Surface Indexes

type RenderPipe struct {
	ID        uint32
	Locations map[string]int32
}

func NewRenderPipe(id uint32) *RenderPipe {
	rProg := &RenderPipe{
		ID:        id,
		Locations: make(map[string]int32, 10),
	}
	return rProg
}

type RenderIndex int

const (
	Standard2D RenderIndex = iota
	Primitive2D
	Primitive2DVariableColor
	Textured2D
	Textured2DVariableColor
) // Render Pipe Indexes

// RENDER SURFACE
type RenderSurface struct {
	sID  uint32
	tID  uint32
	size Vec2
}

func NewRenderSurface(sID uint32, tID uint32, size Vec2) *RenderSurface {
	return &RenderSurface{
		sID:  sID,
		tID:  tID,
		size: size,
	}
}

type GraphicsInterface interface {
	Init()
	Run(func())
	Teardown()
	GetWindowSize() V.F32Vec2
	AddRenderPipe(rendIndex RenderIndex, vShader *Shader, fShader *Shader)
	AddTexture(texIndex TextureIndex, texture *Texture)
	AddRenderSurface(surfIndex SurfaceIndex, texIndex TextureIndex, size Vec2)
	ClearSurface(baseColor *Color)
	ClearSurfaceArea(surfIndex SurfaceIndex, baseColor *Color, rect Rect2D)

	DrawBatchIndexedTriangles2D()
	AddVertexToBatch(pos Vec2, color *Color, uv Vec2) (index uint16)
	AddIndexesToBatch(indexes ...uint16)
	//DrawPrimitiveVertexArray2D(verts []Vec2, color *Color, mode VertexMode)
	//DrawTexturedVertexArray2D(texIndex TextureIndex, destVerts []Vec2, sourceVerts []Vec2, color *Color, mode VertexMode, blendAlpha bool)
	// Drawing modes
	DrawToScreen(op func())
	DrawToSurface(surfIndex SurfaceIndex, op func())
	//DrawUsingRenderPipe(rendIndex RenderIndex, op func())
}

type InputInterface interface {
	SetClipboardText(text string)
	GetClipboardText() string
	// Mouse Input
	GetMouseButtonState(button MouseButton) InputState
	GetMousePosition() Vec2
	SetCallbackOnMouseWheelScroll(op func(offset Vec2))
	SetCallbackOnMouseMove(op func(pos Vec2))
	SetCallbackOnMouseButton(op func(button MouseButton, state InputState))
	// Keyboard Input
	GetKeyboardKeyState(key KeyboardKey) InputState
	SetCallbackOnRuneInput(op func(r rune))
	SetCallbackOnKeyPress(op func(key KeyboardKey, state InputState, mods KeyboardMod))
	// Touch Input
	//// TODO:
	// Controller Input
	//// TODO:
}

type SystemSolution struct {
	lib   GraphicsInterface
	fonts map[FontIndex]*QuadPolyFont
	lock  *sync.Mutex
}

var App *SystemSolution

func NewSystemSolution(lib GraphicsInterface) *SystemSolution {
	return &SystemSolution{
		lib:  lib,
		lock: &sync.Mutex{},
	}
}

// Lifetime
func (s *SystemSolution) Init() {
	s.lib.Init()
	s.fonts = make(map[FontIndex]*QuadPolyFont)
	s.AddFont(PlaniTechFontSolid, BuildQuadPolyFont(PlaniTechVBuilder, Vec2{20, 34}, 3.5, 0, 8, 18))
	s.AddFont(PlaniTechFontOutline, BuildQuadPolyFont(PlaniTechVBuilder, Vec2{20, 34}, 7, 0, 8, 18))
	s.AddFont(PlaniTechFontShadow, BuildQuadPolyFont(PlaniTechVBuilder, Vec2{20, 34}, 9, 0, 8, 18))
}
func (s *SystemSolution) Run(op func()) {
	s.lib.Run(op)
}
func (s *SystemSolution) Teardown() {
	s.lib.Teardown()
}
func (s *SystemSolution) ObtainLock(op func()) {
	s.lock.Lock()
	op()
	s.lock.Unlock()
}

// Tools
func (s *SystemSolution) SetClipboardText(text string) {
	s.lib.SetClipboardText(text)
}
func (s *SystemSolution) GetClipboardText() string {
	return s.lib.GetClipboardText()
}
func (s *SystemSolution) GetWindowSize() Vec2 {
	return s.lib.GetWindowSize()
}

// Asset Linking
func (s *SystemSolution) AddRenderPipe(pIndex RenderIndex, vShader *Shader, fShader *Shader) {
	s.lib.AddRenderPipe(pIndex, vShader, fShader)
}
func (s *SystemSolution) AddTexture(index TextureIndex, texture *Texture) {
	s.lib.AddTexture(index, texture)
}
func (s *SystemSolution) AddRenderSurface(surfIndex SurfaceIndex, texIndex TextureIndex, size Vec2) {
	s.lib.AddRenderSurface(surfIndex, texIndex, size)
}
func (s *SystemSolution) AddFont(fontIndex FontIndex, font *QuadPolyFont) {
	s.fonts[fontIndex] = font
}
func (s *SystemSolution) GetFont(fontIndex FontIndex) *QuadPolyFont {
	return s.fonts[fontIndex]
}

// Draw Modes
func (s *SystemSolution) DrawToScreen(op func()) {
	s.lib.DrawToScreen(op)
}
func (s *SystemSolution) DrawToSurface(surfIndex SurfaceIndex, op func()) {
	s.lib.DrawToSurface(surfIndex, op)
}

//func (s *SystemSolution) DrawUsingRenderPipe(rendIndex RenderIndex, op func()) {
//	s.lib.DrawUsingRenderPipe(rendIndex, op)
//}
// Basic Draw Functions
func (s *SystemSolution) ClearSurface(baseColor *Color) {
	s.lib.ClearSurface(baseColor)
}
func (s *SystemSolution) ClearSurfaceArea(surfIndex SurfaceIndex, baseColor *Color, rect Rect2D) {
	s.lib.ClearSurfaceArea(surfIndex, baseColor, rect)
}
func (s *SystemSolution) DrawBatchIndexedTriangles2D() {
	s.lib.DrawBatchIndexedTriangles2D()
}
func (s *SystemSolution) AddVertexToBatch(pos Vec2, color *Color, uv Vec2) (index uint16) {
	return s.lib.AddVertexToBatch(pos, color, uv)
}
func (s *SystemSolution) AddIndexesToBatch(indexes ...uint16) {
	s.lib.AddIndexesToBatch(indexes...)
}

//func (s *SystemSolution) DrawPrimitiveVertexArray2D(verts []Vec2, color *Color, mode VertexMode) {
//	s.lib.DrawPrimitiveVertexArray2D(verts, color, mode)
//}
//func (s *SystemSolution) DrawTexturedVertexArray2D(texIndex TextureIndex, destVerts []Vec2, sourceVerts []Vec2, color *Color, mode VertexMode, blendAlpha bool) {
//	s.lib.DrawTexturedVertexArray2D(texIndex, destVerts, sourceVerts, color, mode, blendAlpha)
//}
// Input Events
func (s *SystemSolution) GetMouseButtonState(button MouseButton) InputState {
	return s.lib.GetMouseButtonState(button)
}
func (s *SystemSolution) SetCallbackOnMouseWheelScroll(op func(offset Vec2)) {
	s.lib.SetCallbackOnMouseWheelScroll(op)
}
func (s *SystemSolution) GetMousePosition() Vec2 {
	return s.lib.GetMousePosition()
}
func (s *SystemSolution) GetKeyboardKeyState(key KeyboardKey) InputState {
	return s.lib.GetKeyboardKeyState(key)
}
func (s *SystemSolution) SetCallbackOnRuneInput(op func(r rune)) {
	s.lib.SetCallbackOnRuneInput(op)
}
func (s *SystemSolution) SetCallbackOnKeyPress(op func(key KeyboardKey, state InputState, mods KeyboardMod)) {
	s.lib.SetCallbackOnKeyPress(op)
}
func (s *SystemSolution) SetCallbackOnMouseMove(op func(pos Vec2)) {
	s.lib.SetCallbackOnMouseMove(op)
}
func (s *SystemSolution) SetCallbackOnMouseButton(op func(button MouseButton, state InputState)) {
	s.lib.SetCallbackOnMouseButton(op)
}

// Advanced Drawing Functions
//func (s *SystemSolution) DrawPixel2D(pos Vec2, color *Color) {
//	s.DrawPrimitiveVertexArray2D([]Vec2{pos}, color, Pixels)
//}
// Polygons and Circles
func (s *SystemSolution) DrawRegularPolygon(pos Vec2, count float32, radius float32, color *Color, rotation float32) {
	count = FFLoor(count)
	idx := make([]uint16, int(count))
	points := PointsOnCircle(count, radius, pos, rotation)
	cen := s.AddVertexToBatch(pos, color, Vec2{-1, -1})
	for i := range points {
		idx[i] = s.AddVertexToBatch(points[i], color, Vec2{-1, -1})
		if i > 0 {
			s.AddIndexesToBatch(cen, idx[i-1], idx[i])
		}
	}
	s.AddIndexesToBatch(cen, idx[len(idx)-1], idx[0])
}
func (s *SystemSolution) DrawRegularPolygonRing(pos Vec2, count float32, innerRadius float32, outerRadius float32, color *Color, rotation float32) {
	count = FFLoor(count)
	idx := make([]uint16, int(count)*2)
	points := PointsOnRing(count, innerRadius, outerRadius, pos, rotation)
	for i := range points {
		idx[i] = s.AddVertexToBatch(points[i], color, Vec2{-1, -1})
	}
	for i := 0; i <= len(idx)-4; i += 2 {
		s.AddIndexesToBatch(idx[i+0], idx[i+1], idx[i+2], idx[i+1], idx[i+3], idx[i+2])
	}
	s.AddIndexesToBatch(idx[len(idx)-2], idx[len(idx)-1], idx[0], idx[len(idx)-1], idx[1], idx[0])
}
func (s *SystemSolution) DrawCircleAutoPoints(pos Vec2, resolution float32, radius float32, color *Color) {
	count := Circumference(radius) / resolution
	s.DrawRegularPolygon(pos, count, radius, color, 0)
}
func (s *SystemSolution) DrawCircleRingAutoPoints(pos Vec2, resolution float32, innerRadius float32, outerRadius float32, color *Color) {
	count := Circumference(outerRadius) / resolution
	s.DrawRegularPolygonRing(pos, count, innerRadius, outerRadius, color, 0)
}
func (s *SystemSolution) DrawCircle(pos Vec2, radius float32, color *Color) {
	s.DrawCircleAutoPoints(pos, 2, radius, color)
}
func (s *SystemSolution) DrawCircleRing(pos Vec2, innerRadius float32, outerRadius float32, color *Color) {
	s.DrawCircleRingAutoPoints(pos, 2, innerRadius, outerRadius, color)
}

// Rectangles
func (s *SystemSolution) DrawRect(rect Rect2D, color *Color) {
	s.DrawRectRotated(rect, color, 0, Vec2{})
}
func (s *SystemSolution) DrawRectOutline(rect Rect2D, color *Color, thickness float32) {
	s.DrawRectOutlineRotated(rect, color, thickness, 0, Vec2{})
}
func (s *SystemSolution) DrawRectRotated(rect Rect2D, color *Color, rotation float32, anchor Vec2) {
	var rectPoints [4]Vec2
	if rotation != 0 {
		rectPoints = rect.RotatedPoints(anchor, rotation)
	} else {
		rectPoints = rect.Points()
	}
	tl := s.AddVertexToBatch(rectPoints[0], color, Vec2{-1, -1})
	tr := s.AddVertexToBatch(rectPoints[1], color, Vec2{-1, -1})
	br := s.AddVertexToBatch(rectPoints[2], color, Vec2{-1, -1})
	bl := s.AddVertexToBatch(rectPoints[3], color, Vec2{-1, -1})
	s.AddIndexesToBatch(bl, tl, br, tl, tr, br)
}
func (s *SystemSolution) DrawRectOutlineRotated(rect Rect2D, color *Color, thickness float32, rotation float32, anchor Vec2) {
	rectOuter := rect.ExpandCopyFromCenter(Vec2{thickness, thickness})
	var rectPointsInner [4]Vec2
	var rectPointsOuter [4]Vec2
	if rotation != 0 {
		rectPointsInner = rect.RotatedPoints(anchor, rotation)
		rectPointsOuter = rectOuter.RotatedPoints(anchor, rotation)
	} else {
		rectPointsInner = rect.Points()
		rectPointsOuter = rectOuter.Points()
	}
	idx := []uint16{
		s.AddVertexToBatch(rectPointsInner[0], color, Vec2{-1, -1}),
		s.AddVertexToBatch(rectPointsOuter[0], color, Vec2{-1, -1}),
		s.AddVertexToBatch(rectPointsInner[1], color, Vec2{-1, -1}),
		s.AddVertexToBatch(rectPointsOuter[1], color, Vec2{-1, -1}),
		s.AddVertexToBatch(rectPointsInner[2], color, Vec2{-1, -1}),
		s.AddVertexToBatch(rectPointsOuter[2], color, Vec2{-1, -1}),
		s.AddVertexToBatch(rectPointsInner[3], color, Vec2{-1, -1}),
		s.AddVertexToBatch(rectPointsOuter[3], color, Vec2{-1, -1}),
	}
	s.AddIndexesToBatch(idx[0], idx[1], idx[2], idx[1], idx[3], idx[2], idx[2], idx[3], idx[4], idx[3], idx[5], idx[4], idx[4], idx[5], idx[6], idx[5], idx[7], idx[6], idx[6], idx[7], idx[0], idx[7], idx[1], idx[0])
}

// Lines
func (s *SystemSolution) DrawLine(a Vec2, b Vec2, thickness float32, color *Color) {
	l := NewLine2D(a, b)
	l1, l2 := l.PerpLines(thickness / 2)
	idx := []uint16{
		s.AddVertexToBatch(l1.A(), color, Vec2{-1, -1}),
		s.AddVertexToBatch(l2.A(), color, Vec2{-1, -1}),
		s.AddVertexToBatch(l1.B(), color, Vec2{-1, -1}),
		s.AddVertexToBatch(l2.B(), color, Vec2{-1, -1}),
	}
	s.AddIndexesToBatch(idx[0], idx[1], idx[2], idx[1], idx[3], idx[2])
}

// Triangle Multi-Strips
func (s *SystemSolution) DrawMultiTriStrips(strips TriStrips, pos Vec2, color *Color) {
	tStrips := strips.Translate(pos)
	s.DrawMultiStripsPreTranslated(tStrips, color)
}
func (s *SystemSolution) DrawMultiStripsPreTranslated(strips TriStrips, color *Color) {
	for _, strip := range strips {
		idx := make([]uint16, len(strip))
		for i := range strip {
			idx[i] = s.AddVertexToBatch(strip[i], color, Vec2{-1, -1})
		}
		for i := 0; i <= len(idx)-4; i += 2 {
			s.AddIndexesToBatch(idx[i+0], idx[i+1], idx[i+2], idx[i+1], idx[i+3], idx[i+2])
		}
	}
}

// Texture Rectangles
func (s *SystemSolution) DrawFromTex(texIndex TextureIndex, source Rect2D, pos Vec2) {
	s.DrawFromTexComplete(texIndex, source, source.WithPos(pos), &ColorWhite, 0, Vec2{}, true)
}
func (s *SystemSolution) DrawFromTexTinted(texIndex TextureIndex, source Rect2D, pos Vec2, color *Color) {
	s.DrawFromTexComplete(texIndex, source, source.WithPos(pos), color, 0, Vec2{}, true)
}
func (s *SystemSolution) DrawFromTexRotated(texIndex TextureIndex, source Rect2D, pos Vec2, rotation float32, anchor Vec2) {
	s.DrawFromTexComplete(texIndex, source, source.WithPos(pos), &ColorWhite, rotation, anchor, true)
}
func (s *SystemSolution) DrawFromTexScaled(texIndex TextureIndex, source Rect2D, pos Vec2, scaleX float32, scaleY float32) {
	scaledSize := Vec2{source.W() * scaleX, source.H() * scaleY}
	s.DrawFromTexComplete(texIndex, source, NewRect2D(pos, scaledSize), &ColorWhite, 0, Vec2{}, true)
}
func (s *SystemSolution) DrawFromTexSourceDestRect(texIndex TextureIndex, source Rect2D, dest Rect2D) {
	s.DrawFromTexComplete(texIndex, source, dest, &ColorWhite, 0, Vec2{}, true)
}
func (s *SystemSolution) DrawFromTexSourceDestRectTinted(texIndex TextureIndex, source Rect2D, dest Rect2D, tint *Color) {
	s.DrawFromTexComplete(texIndex, source, dest, tint, 0, Vec2{}, true)
}
func (s *SystemSolution) DrawFromTexTintedRotated(texIndex TextureIndex, source Rect2D, pos Vec2, color *Color, rotation float32, anchor Vec2) {
	s.DrawFromTexComplete(texIndex, source, source.WithPos(pos), color, rotation, anchor, true)
}
func (s *SystemSolution) DrawFromTexTintedScaled(texIndex TextureIndex, source Rect2D, pos Vec2, color *Color, scaleX float32, scaleY float32) {
	scaledSize := Vec2{source.W() * scaleX, source.H() * scaleY}
	s.DrawFromTexComplete(texIndex, source, NewRect2D(pos, scaledSize), color, 0, Vec2{}, true)
}
func (s *SystemSolution) DrawFromTexRotatedScaled(texIndex TextureIndex, source Rect2D, pos Vec2, rotation float32, anchor Vec2, scaleX float32, scaleY float32) {
	scaledSize := Vec2{source.W() * scaleX, source.H() * scaleY}
	s.DrawFromTexComplete(texIndex, source, NewRect2D(pos, scaledSize), &ColorWhite, rotation, anchor, true)
}
func (s *SystemSolution) DrawFromTexTintedRotatedScaled(texIndex TextureIndex, source Rect2D, pos Vec2, color *Color, rotation float32, anchor Vec2, scaleX float32, scaleY float32) {
	scaledSize := Vec2{source.W() * scaleX, source.H() * scaleY}
	s.DrawFromTexComplete(texIndex, source, NewRect2D(pos, scaledSize), color, rotation, anchor, true)
}
func (s *SystemSolution) DrawFromTexComplete(texIndex TextureIndex, source Rect2D, dest Rect2D, color *Color, rotation float32, anchor Vec2, blendAlpha bool) {
	var dPoints [4]Vec2
	if rotation != 0 {
		dPoints = dest.RotatedPoints(anchor, rotation)
	} else {
		dPoints = dest.Points()
	}
	sPoints := source.Points()
	tl := s.AddVertexToBatch(dPoints[0], color, sPoints[0])
	tr := s.AddVertexToBatch(dPoints[1], color, sPoints[1])
	br := s.AddVertexToBatch(dPoints[2], color, sPoints[2])
	bl := s.AddVertexToBatch(dPoints[3], color, sPoints[3])
	s.AddIndexesToBatch(bl, tl, br, tl, tr, br)
}

// Vector Text
func (s *SystemSolution) DrawQuadVecText(fontIndex FontIndex, text string, pos Vec2, color *Color, textSize float32) {
	font := s.fonts[fontIndex]
	x, y := pos.X(), pos.Y()
	ratio := textSize / font.scale.Y()
	runes := []rune(text)
	for idx, c := range runes {
		if c == ' ' {
			x += font.scale.W() * ratio
			continue
		}
		if c == '\n' {
			x = pos.X()
			y += (font.scale.Y() + font.lineSpacing) * ratio
			continue
		}
		char, exists := font.glyphs[c]
		if !exists {
			char, exists = font.glyphs['�']
			if !exists {
				s.DrawRect(NewRect2D(Vec2{x, y}, font.scale.Mag(ratio)), color)
				x += font.scale.Mag(ratio).W() + (font.charSpacing * ratio)
				continue
			}
		}
		var cStrips TriStrips
		if (c == '"' || c == '\'') && (idx == 0 || unicode.IsSpace(runes[idx-1])) && (idx+1 == len(runes) || unicode.IsPrint(runes[idx+1])) {
			cStrips = char.StripsFlipX()
		} else {
			cStrips = char.strips
		}
		cStrips = cStrips.Scale(Vec2{ratio, ratio})
		scaledWidth := char.size.W() * ratio
		s.DrawMultiTriStrips(cStrips, Vec2{x, y}, color)
		x += scaledWidth + (font.charSpacing * ratio)
	}
}

// Sprite Instance
func (s *SystemSolution) DrawSpriteInstanceTinted(sInst *SpriteInstance, pos Vec2, color *Color) {
	frame := sInst.GetFrame()
	source := frame.texRect
	destPos := frame.drawOffset.Add(pos)
	dest := NewRect2D(destPos, source.Size())
	s.DrawFromTexComplete(frame.texIndex, source, dest, color, 0, Vec2{}, true)
}
func (s *SystemSolution) DrawSpriteInstanceDestRectTinted(sInst *SpriteInstance, dest Rect2D, color *Color) {
	frame := sInst.GetFrame()
	source := frame.texRect
	scale := dest.Size().Div(source.Size())
	destFinal := dest.TranslateCopy(frame.drawOffset.Mult(scale))
	s.DrawFromTexComplete(frame.texIndex, source, destFinal, color, 0, Vec2{}, true)
}
