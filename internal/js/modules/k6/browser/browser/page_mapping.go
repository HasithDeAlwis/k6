package browser

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/grafana/sobek"

	"go.k6.io/k6/internal/js/modules/k6/browser/common"
	"go.k6.io/k6/internal/js/modules/k6/browser/k6ext"
)

// mapPage to the JS module.
//
//nolint:funlen
func mapPage(vu moduleVU, p *common.Page) mapping { //nolint:gocognit,cyclop
	rt := vu.Runtime()
	maps := mapping{
		"bringToFront": func() *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.BringToFront() //nolint:wrapcheck
			})
		},
		"check": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameCheckOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing new frame check options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.Check(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"click": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts, err := parseFrameClickOptions(vu.Context(), opts, p.MainFrame().Timeout())
			if err != nil {
				return nil, err
			}

			return k6ext.Promise(vu.Context(), func() (any, error) {
				err := p.Click(selector, popts)
				return nil, err //nolint:wrapcheck
			}), nil
		},
		"close": func(opts sobek.Value) *sobek.Promise {
			// TODO when opts are implemented for this function pares them here before calling k6ext.Promise and doing it
			// in a goroutine off the event loop. As that will race with anything running on the event loop.
			return k6ext.Promise(vu.Context(), func() (any, error) {
				// It's safe to close the taskqueue for this targetID (if one
				// exists).
				vu.close(p.TargetID())

				return nil, p.Close() //nolint:wrapcheck
			})
		},
		"content": func() *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.Content() //nolint:wrapcheck
			})
		},
		"context": func() mapping {
			return mapBrowserContext(vu, p.Context())
		},
		"dblclick": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameDblClickOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing double click options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.Dblclick(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"dispatchEvent": func(selector, typ string, eventInit, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameDispatchEventOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing page dispatch event options: %w", err)
			}
			earg := exportArg(eventInit)
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.DispatchEvent(selector, typ, earg, popts) //nolint:wrapcheck
			}), nil
		},
		"emulateMedia": func(opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewPageEmulateMediaOptions(p)
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing emulateMedia options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.EmulateMedia(popts) //nolint:wrapcheck
			}), nil
		},
		"emulateVisionDeficiency": func(typ string) *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.EmulateVisionDeficiency(typ) //nolint:wrapcheck
			})
		},
		"evaluate": func(pageFunc sobek.Value, gargs ...sobek.Value) (*sobek.Promise, error) {
			if sobekEmptyString(pageFunc) {
				return nil, fmt.Errorf("evaluate requires a page function")
			}
			funcString := pageFunc.String()
			gopts := exportArgs(gargs)
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.Evaluate(funcString, gopts...)
			}), nil
		},
		"evaluateHandle": func(pageFunc sobek.Value, gargs ...sobek.Value) (*sobek.Promise, error) {
			if sobekEmptyString(pageFunc) {
				return nil, fmt.Errorf("evaluateHandle requires a page function")
			}
			funcString := pageFunc.String()
			gopts := exportArgs(gargs)
			return k6ext.Promise(vu.Context(), func() (any, error) {
				jsh, err := p.EvaluateHandle(funcString, gopts...)
				if err != nil {
					return nil, err //nolint:wrapcheck
				}
				return mapJSHandle(vu, jsh), nil
			}), nil
		},
		"fill": func(selector string, value string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameFillOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing fill options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.Fill(selector, value, popts) //nolint:wrapcheck
			}), nil
		},
		"focus": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameBaseOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing focus options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.Focus(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"frames": func() *sobek.Object {
			var (
				mfrs []mapping
				frs  = p.Frames()
			)
			for _, fr := range frs {
				mfrs = append(mfrs, mapFrame(vu, fr))
			}
			return rt.ToValue(mfrs).ToObject(rt)
		},
		"getAttribute": func(selector string, name string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameBaseOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing getAttribute options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				s, ok, err := p.GetAttribute(selector, name, popts)
				if err != nil {
					return nil, err //nolint:wrapcheck
				}
				if !ok {
					return nil, nil //nolint:nilnil
				}
				return s, nil
			}), nil
		},
		"getByRole": func(role string, opts sobek.Value) (*sobek.Object, error) {
			popts := parseGetByRoleOptions(vu.Context(), opts)

			ml := mapLocator(vu, p.GetByRole(role, popts))
			return rt.ToValue(ml).ToObject(rt), nil
		},
		"getByAltText": func(alt sobek.Value, opts sobek.Value) (*sobek.Object, error) {
			palt, popts := parseGetByBaseOptions(vu.Context(), alt, false, opts)

			ml := mapLocator(vu, p.GetByAltText(palt, popts))
			return rt.ToValue(ml).ToObject(rt), nil
		},
		"getByLabel": func(label sobek.Value, opts sobek.Value) (*sobek.Object, error) {
			plabel, popts := parseGetByBaseOptions(vu.Context(), label, true, opts)

			ml := mapLocator(vu, p.GetByLabel(plabel, popts))
			return rt.ToValue(ml).ToObject(rt), nil
		},
		"getByPlaceholder": func(placeholder sobek.Value, opts sobek.Value) (*sobek.Object, error) {
			pplaceholder, popts := parseGetByBaseOptions(vu.Context(), placeholder, false, opts)

			ml := mapLocator(vu, p.GetByPlaceholder(pplaceholder, popts))
			return rt.ToValue(ml).ToObject(rt), nil
		},
		"getByTitle": func(title sobek.Value, opts sobek.Value) (*sobek.Object, error) {
			ptitle, popts := parseGetByBaseOptions(vu.Context(), title, false, opts)

			ml := mapLocator(vu, p.GetByTitle(ptitle, popts))
			return rt.ToValue(ml).ToObject(rt), nil
		},
		"getByTestId": func(testID sobek.Value) (*sobek.Object, error) {
			ptestID := parseStringOrRegex(testID, false)

			ml := mapLocator(vu, p.GetByTestID(ptestID))
			return rt.ToValue(ml).ToObject(rt), nil
		},
		"getByText": func(text sobek.Value, opts sobek.Value) (*sobek.Object, error) {
			ptext, popts := parseGetByBaseOptions(vu.Context(), text, true, opts)

			ml := mapLocator(vu, p.GetByText(ptext, popts))
			return rt.ToValue(ml).ToObject(rt), nil
		},
		"goto": func(url string, opts sobek.Value) (*sobek.Promise, error) {
			gopts := common.NewFrameGotoOptions(
				p.Referrer(),
				p.NavigationTimeout(),
			)
			if err := gopts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing page navigation options to %q: %w", url, err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				resp, err := p.Goto(url, gopts)
				if err != nil {
					return nil, err //nolint:wrapcheck
				}

				return mapResponse(vu, resp), nil
			}), nil
		},
		"hover": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameHoverOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing hover options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.Hover(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"innerHTML": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameInnerHTMLOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing inner HTML options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.InnerHTML(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"innerText": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameInnerTextOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing inner text options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.InnerText(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"inputValue": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameInputValueOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing input value options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.InputValue(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"isChecked": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameIsCheckedOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parse isChecked options of selector %q: %w", selector, err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.IsChecked(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"isClosed": p.IsClosed,
		"isDisabled": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameIsDisabledOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parse isDisabled options of selector %q: %w", selector, err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.IsDisabled(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"isEditable": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameIsEditableOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parse isEditabled options of selector %q: %w", selector, err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.IsEditable(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"isEnabled": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameIsEnabledOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parse isEnabled options of selector %q: %w", selector, err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.IsEnabled(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"isHidden": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameIsHiddenOptions()
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parse isHidden options of selector %q: %w", selector, err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.IsHidden(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"isVisible": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameIsVisibleOptions()
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing isVisible options of selector %q: %w", selector, err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.IsVisible(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"keyboard": mapKeyboard(vu, p.GetKeyboard()),
		"locator": func(selector string, opts sobek.Value) *sobek.Object {
			ml := mapLocator(vu, p.Locator(selector, opts))
			return rt.ToValue(ml).ToObject(rt)
		},
		"mainFrame": func() *sobek.Object {
			mf := mapFrame(vu, p.MainFrame())
			return rt.ToValue(mf).ToObject(rt)
		},
		"mouse": mapMouse(vu, p.GetMouse()),
		"on":    mapPageOn(vu, p),
		"opener": func() *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.Opener(), nil
			})
		},
		"press": func(selector string, key string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFramePressOptions(p.Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing press options of selector %q: %w", selector, err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.Press(selector, key, popts) //nolint:wrapcheck
			}), nil
		},
		"reload": func(opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewPageReloadOptions(common.LifecycleEventLoad, p.NavigationTimeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing reload options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				resp, err := p.Reload(popts)
				if err != nil {
					return nil, err //nolint:wrapcheck
				}
				if resp == nil {
					return nil, nil //nolint:nilnil
				}
				return mapResponse(vu, resp), nil
			}), nil
		},
		"route": mapPageRoute(vu, p),
		"screenshot": func(opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewPageScreenshotOptions()
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing page screenshot options: %w", err)
			}

			rt := vu.Runtime()
			promise, res, rej := rt.NewPromise()
			callback := vu.RegisterCallback()
			go func() {
				bb, err := p.Screenshot(popts, vu.filePersister)
				if err != nil {
					callback(func() error {
						return rej(err)
					})
					return
				}

				callback(func() error {
					return res(rt.NewArrayBuffer(bb))
				})
			}()

			return promise, nil
		},
		"selectOption": func(selector string, values sobek.Value, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameSelectOptionOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing select option options: %w", err)
			}

			convValues, err := common.ConvertSelectOptionValues(vu.Runtime(), values)
			if err != nil {
				return nil, fmt.Errorf("parsing select options values: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.SelectOption(selector, convValues, popts) //nolint:wrapcheck
			}), nil
		},
		"setChecked": func(selector string, checked bool, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameCheckOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing frame set check options: %w", err)
			}

			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.SetChecked(selector, checked, popts) //nolint:wrapcheck
			}), nil
		},
		"setContent": func(html string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameSetContentOptions(p.MainFrame().NavigationTimeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing setContent options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.SetContent(html, popts) //nolint:wrapcheck
			}), nil
		},
		"setDefaultNavigationTimeout": p.SetDefaultNavigationTimeout,
		"setDefaultTimeout":           p.SetDefaultTimeout,
		"setExtraHTTPHeaders": func(headers map[string]string) *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.SetExtraHTTPHeaders(headers) //nolint:wrapcheck
			})
		},
		"setInputFiles": func(selector string, files sobek.Value, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameSetInputFilesOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing setInputFiles options: %w", err)
			}

			pfiles := new(common.Files)
			if err := pfiles.Parse(vu.Context(), files); err != nil {
				return nil, fmt.Errorf("parsing setInputFiles parameter: %w", err)
			}

			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.SetInputFiles(selector, pfiles, popts) //nolint:wrapcheck
			}), nil
		},
		"setViewportSize": func(viewportSize sobek.Value) (*sobek.Promise, error) {
			s := new(common.Size)
			if err := s.Parse(vu.Context(), viewportSize); err != nil {
				return nil, fmt.Errorf("parsing viewport size: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.SetViewportSize(s) //nolint:wrapcheck
			}), nil
		},
		"tap": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameTapOptions(p.Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing page tap options: %w", err)
			}
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.Tap(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"textContent": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameTextContentOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing text content options: %w", err)
			}

			return k6ext.Promise(vu.Context(), func() (any, error) {
				s, ok, err := p.TextContent(selector, popts)
				if err != nil {
					return nil, err //nolint:wrapcheck
				}
				if !ok {
					return nil, nil //nolint:nilnil
				}
				return s, nil
			}), nil
		},
		"throttleCPU": func(cpuProfile common.CPUProfile) *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.ThrottleCPU(cpuProfile) //nolint:wrapcheck
			})
		},
		"throttleNetwork": func(networkProfile common.NetworkProfile) *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.ThrottleNetwork(networkProfile) //nolint:wrapcheck
			})
		},
		"title": func() *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				return p.Title() //nolint:wrapcheck
			})
		},
		"touchscreen": mapTouchscreen(vu, p.GetTouchscreen()),
		"type": func(selector string, text string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameTypeOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing type options: %w", err)
			}

			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.Type(selector, text, popts) //nolint:wrapcheck
			}), nil
		},
		"uncheck": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameUncheckOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing frame uncheck options %q: %w", selector, err)
			}

			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.Uncheck(selector, popts) //nolint:wrapcheck
			}), nil
		},
		"url":          p.URL,
		"viewportSize": p.ViewportSize,
		"waitForFunction": func(pageFunc, opts sobek.Value, args ...sobek.Value) (*sobek.Promise, error) {
			js, popts, pargs, err := parseWaitForFunctionArgs(
				vu.Context(), p.Timeout(), pageFunc, opts, args...,
			)
			if err != nil {
				return nil, fmt.Errorf("page waitForFunction: %w", err)
			}

			return k6ext.Promise(vu.Context(), func() (result any, reason error) {
				return p.WaitForFunction(js, popts, pargs...) //nolint:wrapcheck
			}), nil
		},
		"waitForLoadState": func(state string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameWaitForLoadStateOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing waitForLoadState %q options: %w", state, err)
			}

			return k6ext.Promise(vu.Context(), func() (any, error) {
				return nil, p.WaitForLoadState(state, popts) //nolint:wrapcheck
			}), nil
		},
		"waitForNavigation": func(opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameWaitForNavigationOptions(p.Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing page wait for navigation options: %w", err)
			}

			// Inject JS regex checker for URL regex pattern matching
			ctx := vu.Context()
			jsRegexChecker, err := injectRegexMatcherScript(ctx, vu, p.TargetID())
			if err != nil {
				return nil, err
			}

			return k6ext.Promise(ctx, func() (result any, reason error) {
				resp, err := p.WaitForNavigation(popts, jsRegexChecker)
				if err != nil {
					return nil, err //nolint:wrapcheck
				}
				return mapResponse(vu, resp), nil
			}), nil
		},
		"waitForSelector": func(selector string, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameWaitForSelectorOptions(p.MainFrame().Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing wait for selector %q options: %w", selector, err)
			}

			return k6ext.Promise(vu.Context(), func() (any, error) {
				eh, err := p.WaitForSelector(selector, popts)
				if err != nil {
					return nil, err //nolint:wrapcheck
				}
				return mapElementHandle(vu, eh), nil
			}), nil
		},
		"waitForTimeout": func(timeout int64) *sobek.Promise {
			return k6ext.Promise(vu.Context(), func() (any, error) {
				p.WaitForTimeout(timeout)
				return nil, nil
			})
		},
		"waitForURL": func(url sobek.Value, opts sobek.Value) (*sobek.Promise, error) {
			popts := common.NewFrameWaitForURLOptions(p.Timeout())
			if err := popts.Parse(vu.Context(), opts); err != nil {
				return nil, fmt.Errorf("parsing waitForURL options: %w", err)
			}

			var val string
			switch url.ExportType() {
			case reflect.TypeOf(string("")):
				val = fmt.Sprintf("'%s'", url.String()) // Strings require quotes
			default: // JS Regex, CSS, numbers or booleans
				val = url.String() // No quotes
			}

			// Inject JS regex checker for URL pattern matching
			ctx := vu.Context()
			jsRegexChecker, err := injectRegexMatcherScript(ctx, vu, p.TargetID())
			if err != nil {
				return nil, err
			}

			return k6ext.Promise(ctx, func() (result any, reason error) {
				return nil, p.WaitForURL(val, popts, jsRegexChecker)
			}), nil
		},
		"workers": func() *sobek.Object {
			var mws []mapping
			for _, w := range p.Workers() {
				mw := mapWorker(vu, w)
				mws = append(mws, mw)
			}
			return rt.ToValue(mws).ToObject(rt)
		},
	}
	maps["$"] = func(selector string) *sobek.Promise {
		return k6ext.Promise(vu.Context(), func() (any, error) {
			eh, err := p.Query(selector)
			if err != nil {
				return nil, err //nolint:wrapcheck
			}
			// ElementHandle can be null when the selector does not match any elements.
			// We do not want to map nil elementHandles since the expectation is a
			// null result in the test script for this case.
			if eh == nil {
				return nil, nil
			}
			ehm := mapElementHandle(vu, eh)

			return ehm, nil
		})
	}
	maps["$$"] = func(selector string) *sobek.Promise {
		return k6ext.Promise(vu.Context(), func() (any, error) {
			ehs, err := p.QueryAll(selector)
			if err != nil {
				return nil, err //nolint:wrapcheck
			}
			var mehs []mapping
			for _, eh := range ehs {
				ehm := mapElementHandle(vu, eh)
				mehs = append(mehs, ehm)
			}
			return mehs, nil
		})
	}

	return maps
}

// mapPageOn maps the requested page.on event to the Sobek runtime.
// It generalizes the handling of page.on events.
func mapPageOn(vu moduleVU, p *common.Page) func(common.PageOnEventName, sobek.Callable) error {
	return func(eventName common.PageOnEventName, handleEvent sobek.Callable) error {
		rt := vu.Runtime()

		pageOnEvents := map[common.PageOnEventName]struct {
			mapp func(vu moduleVU, event common.PageOnEvent) mapping
			init func() error // If set, runs before the event handler.
			wait bool         // Whether to wait for the handler to complete.
		}{
			common.EventPageConsoleAPICalled: {
				mapp: mapConsoleMessage,
				wait: false,
			},
			common.EventPageMetricCalled: {
				mapp: mapMetricEvent,
				init: prepK6BrowserRegExChecker(rt),
				wait: true,
			},
			common.EventPageRequestCalled: {
				mapp: mapRequestEvent,
				wait: false,
			},
			common.EventPageResponseCalled: {
				mapp: mapResponseEvent,
				wait: false,
			},
		}
		pageOnEvent, ok := pageOnEvents[eventName]
		if !ok {
			return fmt.Errorf("unknown page on event: %q", eventName)
		}

		// Initializes the environment for the event handler if necessary.
		if pageOnEvent.init != nil {
			if err := pageOnEvent.init(); err != nil {
				return fmt.Errorf("initiating page.on('%s'): %w", eventName, err)
			}
		}

		ctx := vu.Context()

		// Run the event handler in the task queue to
		// ensure that the handler is executed on the event loop.
		tq := vu.get(ctx, p.TargetID())
		eventHandler := func(event common.PageOnEvent) error {
			mapping := pageOnEvent.mapp(vu, event)

			done := make(chan struct{})

			tq.Queue(func() error {
				defer close(done)

				_, err := handleEvent(
					sobek.Undefined(),
					rt.ToValue(mapping),
				)
				if err != nil {
					return fmt.Errorf("executing page.on('%s') handler: %w", eventName, err)
				}

				return nil
			})

			if pageOnEvent.wait {
				select {
				case <-done:
				case <-ctx.Done():
					return errors.New("iteration ended before page.on handler completed executing")
				}
			}

			return nil
		}

		return p.On(eventName, eventHandler) //nolint:wrapcheck
	}
}

// prepK6BrowserRegExChecker is a helper function to check the regex pattern
// on Sobek runtime. Unlike Go's regexp package, Sobek's runtime checks
// regex patterns using JavaScript's regular expression features.
func prepK6BrowserRegExChecker(rt *sobek.Runtime) func() error {
	return func() error {
		_, err := rt.RunString(`
			function _k6BrowserCheckRegEx(pattern, url) {
				return pattern.test(url);
			}
		`)
		if err != nil {
			return fmt.Errorf("evaluating regex function: %w", err)
		}

		return nil
	}
}

// injectRegexMatcherScript injects a JavaScript regex checker function into the runtime
// for URL pattern matching. This handles regex patterns only using JavaScript's regex
// engine for consistency. It returns a function that can be used to check if a URL
// matches a given pattern in the JS runtime's eventloop.
//
// Do not call this off the main thread (not even from within a promise). The returned
// JSRegexChecker can be called from off the main thread (i.e. in a new goroutine) since
// it will queue up the checker on the event loop.
func injectRegexMatcherScript(ctx context.Context, vu moduleVU, targetID string) (common.JSRegexChecker, error) {
	rt := vu.Runtime()

	err := prepK6BrowserRegExChecker(rt)()
	if err != nil {
		return nil, fmt.Errorf("preparing k6 browser regex checker: %w", err)
	}

	tq := vu.get(ctx, targetID)

	return func(pattern, url string) (bool, error) {
		var (
			result bool
			err    error
		)

		done := make(chan struct{})

		tq.Queue(func() error {
			defer close(done)

			// Regex pattern is unquoted string whereas the url needs to be quoted
			// so that it is treated as a string.
			js := fmt.Sprintf(`_k6BrowserCheckRegEx(%s, '%s')`, pattern, url)

			val, jsErr := rt.RunString(js)
			if jsErr != nil {
				err = fmt.Errorf("evaluating pattern: %w", jsErr)
				return nil
			}

			result = val.ToBoolean()
			return nil
		})

		select {
		case <-done:
		case <-ctx.Done():
			err = fmt.Errorf("context canceled while evaluating URL pattern")
		}

		return result, err
	}, nil
}

func parseWaitForFunctionArgs(
	ctx context.Context, timeout time.Duration, pageFunc, opts sobek.Value, gargs ...sobek.Value,
) (string, *common.FrameWaitForFunctionOptions, []any, error) {
	popts := common.NewFrameWaitForFunctionOptions(timeout)
	err := popts.Parse(ctx, opts)
	if err != nil {
		return "", nil, nil, fmt.Errorf("parsing waitForFunction options: %w", err)
	}

	js := pageFunc.ToString().String()
	_, isCallable := sobek.AssertFunction(pageFunc)
	if !isCallable {
		js = fmt.Sprintf("() => (%s)", js)
	}

	return js, popts, exportArgs(gargs), nil
}

// parseStringOrRegex parses a sobek.Value to return either a quoted string if it was a string,
// or a raw string if it was a JS RegExp object or another type.
//
// Some getBy* APIs work with single quotes and some work with double quotes.
// This inconsistency seems to stem from the injected code copied from
// Playwright itself.
//
// I would prefer not to change the copied injected script code from Playwright
// so that it is easier to copy over updates/fixes from Playwright when we need
// to.
func parseStringOrRegex(v sobek.Value, doubleQuote bool) string {
	var a string
	switch v.ExportType() {
	case reflect.TypeOf(string("")): // text values require quotes
		if doubleQuote {
			a = fmt.Sprintf(`"%s"`, v.String())
		} else {
			a = fmt.Sprintf("'%s'", v.String())
		}
	case reflect.TypeOf(map[string]interface{}(nil)): // JS RegExp
		a = v.String() // No quotes
	default: // CSS, numbers or booleans
		a = v.String() // No quotes
	}
	return a
}

// parseGetByRoleOptions parses the GetByRole options from the Sobek.Value.
func parseGetByRoleOptions(ctx context.Context, opts sobek.Value) *common.GetByRoleOptions {
	if !sobekValueExists(opts) {
		return nil
	}

	o := &common.GetByRoleOptions{}

	rt := k6ext.Runtime(ctx)

	obj := opts.ToObject(rt)
	for _, k := range obj.Keys() {
		switch k {
		case "checked":
			val := obj.Get(k).ToBoolean()
			o.Checked = &val
		case "disabled":
			val := obj.Get(k).ToBoolean()
			o.Disabled = &val
		case "exact":
			val := obj.Get(k).ToBoolean()
			o.Exact = &val
		case "expanded":
			val := obj.Get(k).ToBoolean()
			o.Expanded = &val
		case "includeHidden":
			val := obj.Get(k).ToBoolean()
			o.IncludeHidden = &val
		case "level":
			val := obj.Get(k).ToInteger()
			o.Level = &val
		case "name":
			val := parseStringOrRegex(obj.Get(k), false)
			o.Name = &val
		case "pressed":
			val := obj.Get(k).ToBoolean()
			o.Pressed = &val
		case "selected":
			val := obj.Get(k).ToBoolean()
			o.Selected = &val
		}
	}

	return o
}

// parseGetByBaseOptions parses the options for the GetBy* APIs and the input
// text/regex.
func parseGetByBaseOptions(
	ctx context.Context,
	input sobek.Value,
	doubleQuote bool,
	opts sobek.Value,
) (string, *common.GetByBaseOptions) {
	a := parseStringOrRegex(input, doubleQuote)

	if !sobekValueExists(opts) {
		return a, nil
	}

	o := &common.GetByBaseOptions{}

	rt := k6ext.Runtime(ctx)

	obj := opts.ToObject(rt)
	for _, k := range obj.Keys() {
		if k == "exact" {
			val := obj.Get(k).ToBoolean()
			o.Exact = &val
		}
	}

	return a, o
}

// mapPageRoute maps the requested page.route event to the Sobek runtime.
func mapPageRoute(vu moduleVU, p *common.Page) func(path sobek.Value, handler sobek.Callable) (*sobek.Promise, error) {
	return func(path sobek.Value, handler sobek.Callable) (*sobek.Promise, error) {
		ctx := vu.Context()
		targetID := p.TargetID()

		// Inject JS regex checker for URL regex pattern matching
		jsRegexChecker, err := injectRegexMatcherScript(ctx, vu, targetID)
		if err != nil {
			return nil, err
		}
		pathStr := parseStringOrRegex(path, false)

		// Run the event handler in the task queue to
		// ensure that the handler is executed on the event loop.
		tq := vu.get(ctx, targetID)
		routeHandler := func(route *common.Route) error {
			done := make(chan error, 1)
			tq.Queue(func() error {
				defer close(done)

				mr := mapRoute(vu, route)
				_, err = handler(
					sobek.Undefined(),
					vu.Runtime().ToValue(mr),
				)
				if err != nil {
					done <- fmt.Errorf("executing page.route('%s') handler: %w", path, err)
					return nil
				}

				return nil
			})

			select {
			case err := <-done:
				return err
			case <-ctx.Done():
				return errors.New("iteration ended before route completed")
			}
		}

		return k6ext.Promise(vu.Context(), func() (any, error) {
			return nil, p.Route(pathStr, routeHandler, jsRegexChecker)
		}), nil
	}
}
